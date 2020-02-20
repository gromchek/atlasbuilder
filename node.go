package main

//import "fmt"

type Node struct {
	X int
	Y int
	W int
	H int

	used bool

	down  *Node
	right *Node
}

type Tree struct {
	Root    *Node
	Padding int
}

func MakeNode(x, y, w, h int) *Node {
	return &Node{x, y, w, h, false, nil, nil}
}

func (n *Node) SplitNode(x, y int) *Node {
	n.used = true
	n.down = MakeNode(n.X, n.Y+y, n.W, n.H-y)
	n.right = MakeNode(n.X+x, n.Y, n.W-x, y)

	return n
}

func (root *Node) FindNode(x, y int) *Node {
	if root.used {
		right := root.right.FindNode(x, y)
		down := root.down.FindNode(x, y)

		if right != nil {
			return right
		}

		if down != nil {
			return down
		}
	} else {
		if x <= root.W && y <= root.H {
			return root
		}
	}

	return nil
}

func (tree *Tree) GrowNode(x, y int) *Node {
	canGrowDown := (x <= tree.Root.W)
	canGrowRight := (y <= tree.Root.H)

	shouldGrowRight := canGrowRight && (tree.Root.H >= (tree.Root.W + x))
	shouldGrowDown := canGrowDown && (tree.Root.W >= (tree.Root.H + y))

	if shouldGrowRight {
		return growRight(tree, x, y)
	} else if shouldGrowDown {
		return growDown(tree, x, y)
	} else if canGrowRight {
		return growRight(tree, x, y)
	} else if canGrowDown {
		return growDown(tree, x, y)
	}

	return nil
}

func growRight(tree *Tree, x, y int) *Node {
	tree.Root = &Node{tree.Padding, tree.Padding, tree.Root.W + x, tree.Root.H, true, tree.Root, MakeNode(tree.Root.W, tree.Padding, x, tree.Root.H)}

	if n := tree.Root.FindNode(x, y); n != nil {
		return n.SplitNode(x, y)
	}

	return nil
}

func growDown(tree *Tree, x, y int) *Node {
	tree.Root = &Node{tree.Padding, tree.Padding, tree.Root.W, tree.Root.H + y, true, MakeNode(tree.Padding, tree.Root.H, tree.Root.W, y), tree.Root}

	if n := tree.Root.FindNode(x, y); n != nil {
		return n.SplitNode(x, y)
	}

	return nil
}
