package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

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

	mux         *http.ServeMux
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
		mux:         http.NewServeMux(),
		openapiYAML: yamlSpec,
		openapiJSON: jsonSpec,
	}
	m.routes()
	return m
}

func (m *Module) Handler() http.Handler { return m.mux }

func (m *Module) ListenAddr() string { return m.config.ListenAddr }

func (m *Module) routes() {
	m.mux.HandleFunc("GET /openapi.yaml", m.handleOpenAPIYAML)
	m.mux.HandleFunc("GET /openapi.json", m.handleOpenAPIJSON)
	m.mux.HandleFunc("GET /docs", m.handleDocs)
	m.mux.HandleFunc("GET /docs/", m.handleDocs)

	HandlerWithOptions(m, StdHTTPServerOptions{
		BaseRouter: m.mux,
		ErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		},
	})
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

func (m *Module) handleOpenAPIYAML(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	_, _ = w.Write(m.openapiYAML)
}

func (m *Module) handleOpenAPIJSON(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(m.openapiJSON)
}

func (m *Module) handleDocs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Splatmaker API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
    <style>body { margin: 0; }</style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: '/openapi.yaml',
        dom_id: '#swagger-ui'
      });
    </script>
  </body>
</html>`))
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
	yamlSpec, err := os.ReadFile(path)
	if err != nil {
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
