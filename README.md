# Atlas Builder

This tool packs all png files in an current folder to atlas (and generating json config).

## Build
`go build atlasBuilder.go node.go`

### Args

```
-grow
      grow size (default true)
-ignore string
      ignore files (default "newatlas.png")
-name string
      out atlas name (default "newatlas")
-pad int
      padding between textures (default 1)
-x int
      atlas width size (default 512)
-y int
      atlas height size (default 512)
```

## License

MIT
