package httpapi

import (
	"embed"
	"io/fs"
)

//go:embed swaggerui/index.html
var swaggerUI embed.FS

//go:embed openapi/openapi.yaml
var embeddedOpenAPIYAML []byte

var swaggerUIFS fs.FS = mustSubFS(swaggerUI, "swaggerui")

func mustSubFS(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
