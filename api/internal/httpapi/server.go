package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vinneyto/splatmaker/api/internal/core"
)

type APIServer struct {
	deps Dependencies
}

func NewAPIServer(deps Dependencies) *APIServer { return &APIServer{deps: deps} }

func (s *APIServer) Healthz(c *gin.Context) { c.JSON(http.StatusOK, HealthResponse{Status: "ok"}) }

func (s *APIServer) ListJobs(c *gin.Context, params ListJobsParams) {
	identity, ok := s.authenticate(c); if !ok { return }
	filter := core.JobListFilter{UserID: identity.UserID, Limit: fromIntPtr(params.Limit), Offset: fromIntPtr(params.Offset), Status: toCoreStatusPtr(params.Status)}
	items, err := s.deps.JobService.ListJobs(c.Request.Context(), filter)
	if err != nil { s.writeDomainError(c, err); return }
	resp := make([]JobSummary, 0, len(items))
	for _, j := range items { resp = append(resp, toAPISummary(j)) }
	c.JSON(http.StatusOK, ListJobsResponse{Items: resp})
}

func (s *APIServer) GetJob(c *gin.Context, jobID string) {
	identity, ok := s.authenticate(c); if !ok { return }
	item, err := s.deps.JobService.GetJob(c.Request.Context(), identity.UserID, jobID)
	if err != nil { s.writeDomainError(c, err); return }
	c.JSON(http.StatusOK, toAPIDetails(*item))
}

func (s *APIServer) GetJobResultUrls(c *gin.Context, jobID string, params GetJobResultUrlsParams) {
	identity, ok := s.authenticate(c); if !ok { return }
	items, err := s.deps.JobService.GetJobResultURLs(c.Request.Context(), identity.UserID, jobID, time.Duration(fromIntPtr(params.TtlSeconds))*time.Second)
	if err != nil { s.writeDomainError(c, err); return }
	resp := make([]JobResultURL, 0, len(items))
	for _, u := range items { resp = append(resp, JobResultURL{Key: u.Key, FileName: u.FileName, Url: u.URL, ExpiresAt: u.ExpiresAt.UTC()}) }
	c.JSON(http.StatusOK, JobResultURLsResponse{Items: resp})
}

func (s *APIServer) authenticate(c *gin.Context) (core.UserIdentity, bool) {
	ctx := c.Request.Context()
	req := core.AuthRequest{
		AuthorizationHeader: c.GetHeader("Authorization"),
		OIDCIdentityHeader:  c.GetHeader("X-Amzn-Oidc-Identity"),
		OIDCDataHeader:      c.GetHeader("X-Amzn-Oidc-Data"),
	}
	if s.deps.AuthRequestAdapter != nil {
		ctx, req = s.deps.AuthRequestAdapter.Adapt(ctx, req)
	}
	identity, err := s.deps.AuthService.Authenticate(ctx, req.AuthorizationHeader)
	if err != nil { s.writeDomainError(c, err); return core.UserIdentity{}, false }
	return identity, true
}

func (s *APIServer) writeDomainError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, core.ErrUnauthorized), errors.Is(err, core.ErrInvalidToken), errors.Is(err, core.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrJobNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrInvalidArgument):
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrNotImplemented):
		c.JSON(http.StatusNotImplemented, ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}

func toAPISummary(j core.JobSummary) JobSummary {
	return JobSummary{JobId: j.JobID, Status: JobStatus(j.Status), ProgressPercent: intPtr(j.ProgressPercent), CurrentStep: j.CurrentStep, IdempotencyKey: j.IdempotencyKey, CreatedAt: j.CreatedAt.UTC(), UpdatedAt: j.UpdatedAt.UTC()}
}

func toAPIDetails(j core.JobDetails) JobDetails {
	outputs := make([]OutputFileRef, 0, len(j.OutputFiles))
	for _, o := range j.OutputFiles {
		var size *int
		if o.SizeBytes != nil { x := int(*o.SizeBytes); size = &x }
		outputs = append(outputs, OutputFileRef{Key: o.Key, FileName: o.FileName, SizeBytes: size})
	}
	source := ""
	if j.SourceRef != nil { source = *j.SourceRef }
	return JobDetails{Summary: toAPISummary(j.Summary), Attempt: j.Attempt, SourceRef: source, SimulateFailure: boolPtr(j.SimulateFailure), ErrorMessage: j.ErrorMessage, StartedAt: utcPtr(j.StartedAt), FinishedAt: utcPtr(j.FinishedAt), LastHeartbeatAt: utcPtr(j.LastHeartbeatAt), OutputFiles: outputs}
}

func utcPtr(v *time.Time) *time.Time { if v == nil { return nil }; u := v.UTC(); return &u }
func intPtr(v int) *int { return &v }
func boolPtr(v bool) *bool { return &v }
func fromIntPtr(v *int) int { if v == nil { return 0 }; return *v }
func toCoreStatusPtr(v *JobStatus) *core.JobStatus { if v == nil { return nil }; s := core.JobStatus(*v); return &s }
