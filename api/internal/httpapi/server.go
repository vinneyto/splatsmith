package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type APIServer struct {
	deps Dependencies
}

func NewAPIServer(deps Dependencies) *APIServer {
	return &APIServer{deps: deps}
}

func (s *APIServer) Healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{
		Status: "ok",
		Mode:   s.deps.Mode,
	})
}

func (s *APIServer) ListJobs(w http.ResponseWriter, r *http.Request, params ListJobsParams) {
	identity, ok := s.authenticate(w, r)
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

	items, err := s.deps.JobService.ListJobs(r.Context(), identity.UserID, limit, offset)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}

	respItems := make([]JobSummary, 0, len(items))
	for _, j := range items {
		respItems = append(respItems, toAPISummary(j))
	}

	writeJSON(w, http.StatusOK, ListJobsResponse{Items: respItems})
}

func (s *APIServer) SubmitJob(w http.ResponseWriter, r *http.Request) {
	identity, ok := s.authenticate(w, r)
	if !ok {
		return
	}

	var req SubmitJobJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json body"})
		return
	}

	submit, err := s.deps.JobService.SubmitJob(r.Context(), identity.UserID, core.SubmitJobRequest{
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
		s.writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusAccepted, SubmitJobResponse{
		Created: submit.Created,
		Job:     toAPIDetails(submit.Job),
	})
}

func (s *APIServer) GetJob(w http.ResponseWriter, r *http.Request, jobID string) {
	identity, ok := s.authenticate(w, r)
	if !ok {
		return
	}
	item, err := s.deps.JobService.GetJob(r.Context(), identity.UserID, jobID)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toAPIDetails(*item))
}

func (s *APIServer) CancelJob(w http.ResponseWriter, r *http.Request, jobID string) {
	identity, ok := s.authenticate(w, r)
	if !ok {
		return
	}
	item, err := s.deps.JobService.CancelJob(r.Context(), identity.UserID, jobID)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, JobMutationResponse{Job: toAPIDetails(*item)})
}

func (s *APIServer) RetryJob(w http.ResponseWriter, r *http.Request, jobID string) {
	identity, ok := s.authenticate(w, r)
	if !ok {
		return
	}
	item, err := s.deps.JobService.RetryJob(r.Context(), identity.UserID, jobID)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, JobMutationResponse{Job: toAPIDetails(*item)})
}

func (s *APIServer) GetJobResultUrls(w http.ResponseWriter, r *http.Request, jobID string, params GetJobResultUrlsParams) {
	identity, ok := s.authenticate(w, r)
	if !ok {
		return
	}

	if jobID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "job_id is required"})
		return
	}

	ttlSeconds := int(s.deps.DefaultResultURLTTL.Seconds())
	if params.TtlSeconds != nil && *params.TtlSeconds > 0 {
		ttlSeconds = *params.TtlSeconds
	}

	urls, err := s.deps.JobService.GetJobResultURLs(r.Context(), identity.UserID, jobID, time.Duration(ttlSeconds)*time.Second)
	if err != nil {
		s.writeDomainError(w, err)
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

func (s *APIServer) authenticate(w http.ResponseWriter, r *http.Request) (core.UserIdentity, bool) {
	identity, err := s.deps.AuthService.Authenticate(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		s.writeDomainError(w, err)
		return core.UserIdentity{}, false
	}
	return identity, true
}

func (s *APIServer) writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, core.ErrUnauthorized), errors.Is(err, core.ErrInvalidToken):
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrJobNotFound):
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrIdempotencyKeyRequired), errors.Is(err, core.ErrInvalidArgument):
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrJobNotCancelable), errors.Is(err, core.ErrJobNotRetryable), errors.Is(err, core.ErrConflict):
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
	case errors.Is(err, core.ErrNotImplemented):
		writeJSON(w, http.StatusNotImplemented, ErrorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
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
