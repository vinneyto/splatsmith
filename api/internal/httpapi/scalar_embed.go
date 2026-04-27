package httpapi

import (
	"embed"
	"io/fs"
)

//go:embed scalarui/index.html
var scalarUI embed.FS

//go:embed openapi/openapi.yaml
var embeddedOpenAPIYAML []byte

var scalarUIFS fs.FS = mustSubFS(scalarUI, "scalarui")

func mustSubFS(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
