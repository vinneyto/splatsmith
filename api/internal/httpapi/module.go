package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vinneyto/splatmaker/api/internal/core"
	"github.com/vinneyto/splatmaker/api/internal/core/services"
	"gopkg.in/yaml.v3"
)

type Dependencies struct {
	Mode                string
	AuthService         *services.AuthService
	AuthRequestAdapter  core.AuthRequestAdapter
	JobService          *services.JobService
	DefaultResultURLTTL time.Duration
}

type Module struct {
	config Config
	deps   Dependencies

	engine      *gin.Engine
	apiServer   *APIServer
	openapiYAML []byte
	openapiJSON []byte
}

func NewModule(cfg Config, deps Dependencies) *Module {
	cfg = cfg.withDefaults()
	if deps.DefaultResultURLTTL <= 0 {
		deps.DefaultResultURLTTL = 900 * time.Second
	}

	yamlSpec, jsonSpec := loadOpenAPISpec()
	m := &Module{
		config:      cfg,
		deps:        deps,
		engine:      gin.New(),
		apiServer:   NewAPIServer(deps),
		openapiYAML: yamlSpec,
		openapiJSON: jsonSpec,
	}
	m.engine.Use(gin.Recovery())
	m.engine.Use(corsMiddleware())
	m.routes()
	return m
}

func (m *Module) Handler() http.Handler { return m.engine }

func (m *Module) ListenAddr() string { return m.config.ListenAddr }

func (m *Module) routes() {
	m.engine.GET("/openapi.yaml", m.handleOpenAPIYAML)
	m.engine.GET("/openapi.json", m.handleOpenAPIJSON)
	m.engine.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/scalar/index.html")
	})
	m.engine.GET("/docs/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/scalar/index.html")
	})
	m.engine.GET("/scalar", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/scalar/")
	})
	m.engine.GET("/scalar/*filepath", gin.WrapH(http.StripPrefix("/scalar/", http.FileServer(http.FS(scalarUIFS)))))

	RegisterHandlersWithOptions(m.engine, m.apiServer, GinServerOptions{
		ErrorHandler: func(c *gin.Context, err error, _ int) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		},
	})
}

func (m *Module) handleOpenAPIYAML(c *gin.Context) {
	c.Data(http.StatusOK, "application/yaml; charset=utf-8", m.openapiYAML)
}

func (m *Module) handleOpenAPIJSON(c *gin.Context) {
	c.Data(http.StatusOK, "application/json; charset=utf-8", m.openapiJSON)
}

func loadOpenAPISpec() ([]byte, []byte) {
	yamlSpec := embeddedOpenAPIYAML
	if len(yamlSpec) == 0 {
		fallback := []byte("openapi: 3.0.3\ninfo:\n  title: Splatmaker API\n  version: 0.1.0\npaths: {}\n")
		return fallback, []byte(`{"openapi":"3.0.3","info":{"title":"Splatmaker API","version":"0.1.0"},"paths":{}}`)
	}

	var v any
	if err := yaml.Unmarshal(yamlSpec, &v); err != nil {
		return yamlSpec, []byte(`{"error":"failed to parse openapi.yaml"}`)
	}
	jsonSpec, err := json.Marshal(v)
	if err != nil {
		return yamlSpec, []byte(`{"error":"failed to marshal openapi spec"}`)
	}
	return yamlSpec, jsonSpec
}
