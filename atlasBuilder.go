package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	atlasName string
	sizeX     int
	sizeY     int
	padSize   int
	ignore    string
	growing   bool
)

func init() {
	flag.StringVar(&atlasName, "name", "newatlas", "out atlas name")
	flag.StringVar(&ignore, "ignore", "newatlas.png", "ignore files")
	flag.IntVar(&sizeX, "x", 512, "atlas width size")
	flag.IntVar(&sizeY, "y", 512, "atlas height size")
	flag.IntVar(&padSize, "pad", 1, "padding between textures")
	flag.BoolVar(&growing, "grow", true, "grow size")
	flag.Parse()
}

func stringInArray(str string, array []string) bool {
	for _, s := range array {
		if s == str {
			return true
		}
	}

	return false
}

type Vec2 = image.Point

type ImageInfo struct {
	Name string
	Size Vec2
	Img  image.Image

	Node *Node
}

type ImgInfo []ImageInfo

func (e ImgInfo) Len() int {
	return len(e)
}

func (e ImgInfo) Less(i, j int) bool {
	//return e[i].Size.X > e[j].Size.X && e[i].Size.Y > e[j].Size.Y
	return (e[i].Size.X * e[i].Size.Y) > (e[j].Size.X * e[j].Size.Y)
}

func (e ImgInfo) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type Atlas struct {
	X int
	Y int
	W int
	H int
}

func main() {
	ignoreList := strings.Split(ignore, " ")
	ignoreList = append(ignoreList, atlasName+".png")
	fileNames := []string{}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".png" || stringInArray(f.Name(), ignoreList) {
			continue
		}
		fileNames = append(fileNames, f.Name())
	}

	var images []ImageInfo

	for _, fileName := range fileNames {
		fileHandler, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("Error open %s file\n", fileName)
			continue
		}
		defer fileHandler.Close()

		src, _, err := image.Decode(fileHandler)

		if err != nil {
			fmt.Printf("Error decode '%s' file\n", fileName)
			continue
		}

		fileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		images = append(images, ImageInfo{fileName, Vec2{src.Bounds().Max.X, src.Bounds().Max.Y}, src, nil})
	}

	sort.Sort(ImgInfo(images))

	atlasMap := make(map[string]map[string]Atlas)
	atlasMap[atlasName] = make(map[string]Atlas)

	tree := Tree{MakeNode(padSize, padSize, sizeX, sizeY), padSize}

	for i := 0; i < len(images); i++ {

		if n := tree.Root.FindNode(images[i].Size.X+padSize, images[i].Size.Y+padSize); n != nil {
			images[i].Node = n.SplitNode(images[i].Size.X+padSize, images[i].Size.Y+padSize)
		} else if growing {
			images[i].Node = tree.GrowNode(images[i].Size.X+padSize, images[i].Size.Y+padSize)
		}

	}

	for _, imageInfo := range images {
		if imageInfo.Node != nil {
			atlasMap[atlasName][imageInfo.Name] = Atlas{imageInfo.Node.X, imageInfo.Node.Y, imageInfo.Size.X, imageInfo.Size.Y}
		} else {
			fmt.Printf("Can't fit '%s'\n", imageInfo.Name)
		}
	}

	jsonString, _ := json.Marshal(atlasMap)
	_ = ioutil.WriteFile(atlasName+".json", jsonString, 0644)

	atlasTexture := image.NewRGBA(image.Rect(0, 0, tree.Root.W, tree.Root.H))

	for k, v := range atlasMap[atlasName] {
		for _, imageInfo := range images {
			if k == imageInfo.Name {
				draw.Draw(atlasTexture, image.Rectangle{image.Point{v.X, v.Y}, image.Point{tree.Root.W, tree.Root.H}}, imageInfo.Img, image.Point{}, draw.Src)
			}
		}
	}

	file, _ := os.Create(atlasName + ".png")
	defer file.Close()

	png.Encode(file, atlasTexture)
}
