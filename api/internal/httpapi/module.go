package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/vinneyto/splatra/api/internal/core"
	"github.com/vinneyto/splatra/api/internal/core/services"
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

	mux     *http.ServeMux
	openapi []byte
}

func NewModule(cfg Config, deps Dependencies) *Module {
	cfg = cfg.withDefaults()
	if deps.DefaultResultURLTTL <= 0 {
		deps.DefaultResultURLTTL = 900 * time.Second
	}

	m := &Module{
		config:  cfg,
		deps:    deps,
		mux:     http.NewServeMux(),
		openapi: buildOpenAPISpecJSON(deps.Mode, int(deps.DefaultResultURLTTL.Seconds())),
	}
	m.routes()
	return m
}

func (m *Module) Handler() http.Handler { return m.mux }

func (m *Module) ListenAddr() string { return m.config.ListenAddr }

func (m *Module) routes() {
	m.mux.HandleFunc("GET /healthz", m.handleHealth)
	m.mux.HandleFunc("GET /openapi.json", m.handleOpenAPI)
	m.mux.HandleFunc("GET /docs", m.handleDocs)
	m.mux.HandleFunc("GET /docs/", m.handleDocs)

	m.mux.HandleFunc("GET /v1/jobs", m.handleListJobs)
	m.mux.HandleFunc("GET /v1/jobs/{jobID}/result-urls", m.handleJobResultURLs)
}

func (m *Module) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"mode":   m.deps.Mode,
	})
}

func (m *Module) handleOpenAPI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(m.openapi)
}

func (m *Module) handleDocs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Splatra API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
    <style>body { margin: 0; }</style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: '/openapi.json',
        dom_id: '#swagger-ui'
      });
    </script>
  </body>
</html>`))
}

func (m *Module) handleListJobs(w http.ResponseWriter, r *http.Request) {
	identity, ok := m.authenticate(w, r)
	if !ok {
		return
	}

	limit := parseIntDefault(r.URL.Query().Get("limit"), 20)
	offset := parseIntDefault(r.URL.Query().Get("offset"), 0)
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

	respItems := make([]jobSummaryResponse, 0, len(items))
	for _, j := range items {
		respItems = append(respItems, jobSummaryResponse{
			JobID:     j.JobID,
			Status:    string(j.Status),
			CreatedAt: j.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: j.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": respItems})
}

func (m *Module) handleJobResultURLs(w http.ResponseWriter, r *http.Request) {
	identity, ok := m.authenticate(w, r)
	if !ok {
		return
	}

	jobID := r.PathValue("jobID")
	if jobID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "jobID is required"})
		return
	}

	ttlSeconds := parseIntDefault(r.URL.Query().Get("ttl_seconds"), int(m.deps.DefaultResultURLTTL.Seconds()))
	if ttlSeconds <= 0 {
		ttlSeconds = int(m.deps.DefaultResultURLTTL.Seconds())
	}

	urls, err := m.deps.JobViewer.GetJobResultURLs(r.Context(), identity.UserID, jobID, time.Duration(ttlSeconds)*time.Second)
	if err != nil {
		m.writeDomainError(w, err)
		return
	}

	respItems := make([]resultURLResponse, 0, len(urls))
	for _, u := range urls {
		respItems = append(respItems, resultURLResponse{
			Key:       u.Key,
			FileName:  u.FileName,
			URL:       u.URL,
			ExpiresAt: u.ExpiresAt.UTC().Format(time.RFC3339),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": respItems})
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
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
	case errors.Is(err, core.ErrJobNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, core.ErrNotImplemented):
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

func parseIntDefault(raw string, def int) int {
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return v
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

type jobSummaryResponse struct {
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type resultURLResponse struct {
	Key       string `json:"key"`
	FileName  string `json:"file_name"`
	URL       string `json:"url"`
	ExpiresAt string `json:"expires_at"`
}
