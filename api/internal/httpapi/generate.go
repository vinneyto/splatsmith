package httpapi

//go:generate cp ../../openapi.yaml ./openapi_embedded.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0 -config ../../oapi-codegen.yaml ../../openapi.yaml
