package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"github.com/vinneyto/splatmaker/api/internal/core"
	"github.com/vinneyto/splatmaker/api/internal/core/services"
)

type Dependencies struct {
	Mode                string
	AuthService         *services.AuthService
	JobViewer           *services.JobViewerService
	DefaultResultURLTTL time.Duration
}

type Module struct {
	config Config
	deps   Dependencies

	engine      *gin.Engine
	openapiYAML []byte
	openapiJSON []byte
}

func NewModule(cfg Config, deps Dependencies) *Module {
	cfg = cfg.withDefaults()
	if deps.DefaultResultURLTTL <= 0 {
		deps.DefaultResultURLTTL = 900 * time.Second
	}

	yamlSpec, jsonSpec := loadOpenAPISpec(cfg.OpenAPIPath)
	m := &Module{
		config:      cfg,
		deps:        deps,
		engine:      gin.New(),
		openapiYAML: yamlSpec,
		openapiJSON: jsonSpec,
	}
	m.engine.Use(gin.Recovery())
	m.routes()
	return m
}

func (m *Module) Handler() http.Handler { return m.engine }

func (m *Module) ListenAddr() string { return m.config.ListenAddr }

func (m *Module) routes() {
	m.engine.GET("/openapi.yaml", m.handleOpenAPIYAML)
	m.engine.GET("/openapi.json", m.handleOpenAPIJSON)
	m.engine.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	})
	m.engine.GET("/docs/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	})
	m.engine.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/swagger/")
	})
	m.engine.GET("/swagger/*filepath", gin.WrapH(http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerUIFS)))))

	apiHandler := HandlerWithOptions(m, StdHTTPServerOptions{
		ErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		},
	})
	m.engine.NoRoute(gin.WrapH(apiHandler))
}

func (m *Module) Healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{
		Status: "ok",
		Mode:   m.deps.Mode,
	})
}

func (m *Module) ListJobs(w http.ResponseWriter, r *http.Request, params ListJobsParams) {
	identity, ok := m.authenticate(w, r)
	if !ok {
		return
	}

	limit := 20
	offset := 0
	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	items, err := m.deps.JobViewer.ListJobs(r.Context(), identity.UserID, limit, offset)
	if err != nil {
		m.writeDomainError(w, err)
		return
	}

	respItems := make([]JobSummary, 0, len(items))
	for _, j := range items {
		respItems = append(respItems, JobSummary{
			JobId:     j.JobID,
			Status:    JobSummaryStatus(j.Status),
			CreatedAt: j.CreatedAt.UTC(),
			UpdatedAt: j.UpdatedAt.UTC(),
		})
	}

	writeJSON(w, http.StatusOK, ListJobsResponse{Items: respItems})
}

func (m *Module) GetJobResultUrls(w http.ResponseWriter, r *http.Request, jobID string, params GetJobResultUrlsParams) {
	identity, ok := m.authenticate(w, r)
	if !ok {
		return
	}

	if jobID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "job_id is required"})
		return
	}

	ttlSeconds := int(m.deps.DefaultResultURLTTL.Seconds())
	if params.TtlSeconds != nil && *params.TtlSeconds > 0 {
		ttlSeconds = *params.TtlSeconds
	}

	urls, err := m.deps.JobViewer.GetJobResultURLs(r.Context(), identity.UserID, jobID, time.Duration(ttlSeconds)*time.Second)
	if err != nil {
		m.writeDomainError(w, err)
		return
	}

	respItems := make([]ResultFileURL, 0, len(urls))
	for _, u := range urls {
		respItems = append(respItems, ResultFileURL{
			Key:       u.Key,
			FileName:  u.FileName,
			Url:       u.URL,
			ExpiresAt: u.ExpiresAt.UTC(),
		})
	}
	writeJSON(w, http.StatusOK, JobResultURLsResponse{Items: respItems})
}

func (m *Module) handleOpenAPIYAML(c *gin.Context) {
	c.Data(http.StatusOK, "application/yaml; charset=utf-8", m.openapiYAML)
}

func (m *Module) handleOpenAPIJSON(c *gin.Context) {
	c.Data(http.StatusOK, "application/json; charset=utf-8", m.openapiJSON)
}

func (m *Module) authenticate(w http.ResponseWriter, r *http.Request) (core.UserIdentity, bool) {
	identity, err := m.deps.AuthService.Authenticate(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		m.writeDomainError(w, err)
		return core.UserIdentity{}, false
	}
	return identity, true
}

func (m *Module) writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, core.ErrUnauthorized), errors.Is(err, core.ErrInvalidToken):
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrJobNotFound):
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrNotImplemented):
		writeJSON(w, http.StatusNotImplemented, ErrorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}

func loadOpenAPISpec(path string) ([]byte, []byte) {
	yamlSpec := embeddedOpenAPIYAML
	if len(yamlSpec) == 0 {
		if fromFile, err := os.ReadFile(path); err == nil {
			yamlSpec = fromFile
		}
	}
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

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
