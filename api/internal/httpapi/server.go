package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vinneyto/splatmaker/api/internal/core"
)

type APIServer struct {
	deps Dependencies
}

func NewAPIServer(deps Dependencies) *APIServer {
	return &APIServer{deps: deps}
}

func (s *APIServer) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status: "ok",
		Mode:   s.deps.Mode,
	})
}

func (s *APIServer) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	result, err := s.deps.LoginService.LoginWithPassword(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		s.writeDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: result.Token,
		TokenType:   Bearer,
		User: AuthUser{
			UserId: result.User.UserID,
			Email:  result.User.Email,
		},
	})
}

func (s *APIServer) ListJobs(c *gin.Context, params ListJobsParams) {
	identity, ok := s.authenticate(c)
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

	items, err := s.deps.JobService.ListJobs(c.Request.Context(), identity.UserID, limit, offset)
	if err != nil {
		s.writeDomainError(c, err)
		return
	}

	respItems := make([]JobSummary, 0, len(items))
	for _, j := range items {
		respItems = append(respItems, toAPISummary(j))
	}

	c.JSON(http.StatusOK, ListJobsResponse{Items: respItems})
}

func (s *APIServer) SubmitJob(c *gin.Context) {
	identity, ok := s.authenticate(c)
	if !ok {
		return
	}

	var req SubmitJobJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json body"})
		return
	}

	submit, err := s.deps.JobService.SubmitJob(c.Request.Context(), identity.UserID, core.SubmitJobRequest{
		IdempotencyKey: req.IdempotencyKey,
		Name:           req.Name,
		SourceRef:      req.SourceRef,
		SimulateFailure: func() bool {
			if req.SimulateFailure == nil {
				return false
			}
			return *req.SimulateFailure
		}(),
	})
	if err != nil {
		s.writeDomainError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, SubmitJobResponse{
		Created: submit.Created,
		Job:     toAPIDetails(submit.Job),
	})
}

func (s *APIServer) GetJob(c *gin.Context, jobID string) {
	identity, ok := s.authenticate(c)
	if !ok {
		return
	}
	item, err := s.deps.JobService.GetJob(c.Request.Context(), identity.UserID, jobID)
	if err != nil {
		s.writeDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, toAPIDetails(*item))
}

func (s *APIServer) CancelJob(c *gin.Context, jobID string) {
	identity, ok := s.authenticate(c)
	if !ok {
		return
	}
	item, err := s.deps.JobService.CancelJob(c.Request.Context(), identity.UserID, jobID)
	if err != nil {
		s.writeDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, JobMutationResponse{Job: toAPIDetails(*item)})
}

func (s *APIServer) RetryJob(c *gin.Context, jobID string) {
	identity, ok := s.authenticate(c)
	if !ok {
		return
	}
	item, err := s.deps.JobService.RetryJob(c.Request.Context(), identity.UserID, jobID)
	if err != nil {
		s.writeDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, JobMutationResponse{Job: toAPIDetails(*item)})
}

func (s *APIServer) GetJobResultUrls(c *gin.Context, jobID string, params GetJobResultUrlsParams) {
	identity, ok := s.authenticate(c)
	if !ok {
		return
	}

	if jobID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "job_id is required"})
		return
	}

	ttlSeconds := int(s.deps.DefaultResultURLTTL.Seconds())
	if params.TtlSeconds != nil && *params.TtlSeconds > 0 {
		ttlSeconds = *params.TtlSeconds
	}

	urls, err := s.deps.JobService.GetJobResultURLs(c.Request.Context(), identity.UserID, jobID, time.Duration(ttlSeconds)*time.Second)
	if err != nil {
		s.writeDomainError(c, err)
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
	c.JSON(http.StatusOK, JobResultURLsResponse{Items: respItems})
}

func (s *APIServer) GetStandardPipelineSettings(c *gin.Context) {
	identity, ok := s.authenticate(c)
	if !ok {
		return
	}

	settings, err := s.deps.PipelineSettingsService.GetStandard(c.Request.Context(), identity.UserID)
	if err != nil {
		s.writeDomainError(c, err)
		return
	}

	apiSettings, err := toAPIPipelineSettings(settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	c.JSON(http.StatusOK, PipelineSettingsResponse{Settings: apiSettings})
}

func (s *APIServer) PutStandardPipelineSettings(c *gin.Context) {
	identity, ok := s.authenticate(c)
	if !ok {
		return
	}

	var req PutStandardPipelineSettingsJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json body"})
		return
	}

	coreSettings, err := toCorePipelineSettings(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid pipeline settings payload"})
		return
	}

	saved, err := s.deps.PipelineSettingsService.SaveStandard(c.Request.Context(), identity.UserID, coreSettings)
	if err != nil {
		s.writeDomainError(c, err)
		return
	}

	apiSettings, err := toAPIPipelineSettings(saved)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	c.JSON(http.StatusOK, PipelineSettingsResponse{Settings: apiSettings})
}

func (s *APIServer) authenticate(c *gin.Context) (core.UserIdentity, bool) {
	identity, err := s.deps.AuthService.Authenticate(c.Request.Context(), c.GetHeader("Authorization"))
	if err != nil {
		s.writeDomainError(c, err)
		return core.UserIdentity{}, false
	}
	return identity, true
}

func (s *APIServer) writeDomainError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, core.ErrUnauthorized), errors.Is(err, core.ErrInvalidToken), errors.Is(err, core.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrJobNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrIdempotencyKeyRequired), errors.Is(err, core.ErrInvalidArgument):
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrJobNotCancelable), errors.Is(err, core.ErrJobNotRetryable), errors.Is(err, core.ErrConflict):
		c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrNotImplemented):
		c.JSON(http.StatusNotImplemented, ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}

func toAPISummary(j core.JobSummary) JobSummary {
	return JobSummary{
		JobId:           j.JobID,
		Status:          JobStatus(j.Status),
		ProgressPercent: j.ProgressPercent,
		CurrentStep:     j.CurrentStep,
		IdempotencyKey:  j.IdempotencyKey,
		CreatedAt:       j.CreatedAt.UTC(),
		UpdatedAt:       j.UpdatedAt.UTC(),
	}
}

func toAPIDetails(j core.JobDetails) JobDetails {
	outputs := make([]OutputFileRef, 0, len(j.OutputFiles))
	for _, o := range j.OutputFiles {
		var size *int
		if o.SizeBytes != nil {
			x := int(*o.SizeBytes)
			size = &x
		}
		outputs = append(outputs, OutputFileRef{
			Key:       o.Key,
			FileName:  o.FileName,
			SizeBytes: size,
		})
	}
	return JobDetails{
		Summary:         toAPISummary(j.Summary),
		Attempt:         j.Attempt,
		SourceRef:       j.SourceRef,
		SimulateFailure: j.SimulateFailure,
		ErrorMessage:    j.ErrorMessage,
		StartedAt:       utcPtr(j.StartedAt),
		FinishedAt:      utcPtr(j.FinishedAt),
		LastHeartbeatAt: utcPtr(j.LastHeartbeatAt),
		OutputFiles:     outputs,
	}
}

func utcPtr(v *time.Time) *time.Time {
	if v == nil {
		return nil
	}
	u := v.UTC()
	return &u
}

func toCorePipelineSettings(in PipelineSettings) (core.PipelineSettings, error) {
	b, err := json.Marshal(in)
	if err != nil {
		return core.PipelineSettings{}, err
	}
	var out core.PipelineSettings
	if err := json.Unmarshal(b, &out); err != nil {
		return core.PipelineSettings{}, err
	}
	return out, nil
}

func toAPIPipelineSettings(in core.PipelineSettings) (PipelineSettings, error) {
	b, err := json.Marshal(in)
	if err != nil {
		return PipelineSettings{}, err
	}
	var out PipelineSettings
	if err := json.Unmarshal(b, &out); err != nil {
		return PipelineSettings{}, err
	}
	return out, nil
}
