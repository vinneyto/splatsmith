package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type JobService struct {
	repo       core.JobRepository
	dispatcher core.JobDispatcher
	resolver   core.ResultURLResolver
}

func NewJobService(repo core.JobRepository, dispatcher core.JobDispatcher, resolver core.ResultURLResolver) *JobService {
	return &JobService{repo: repo, dispatcher: dispatcher, resolver: resolver}
}

func (s *JobService) SubmitJob(ctx context.Context, userID string, req core.SubmitJobRequest) (core.SubmitJobResult, error) {
	idempotencyKey := strings.TrimSpace(req.IdempotencyKey)
	if idempotencyKey == "" {
		return core.SubmitJobResult{}, core.ErrIdempotencyKeyRequired
	}
	req.IdempotencyKey = idempotencyKey

	existing, err := s.repo.FindByIdempotencyKey(ctx, userID, idempotencyKey)
	if err != nil {
		return core.SubmitJobResult{}, err
	}
	if existing != nil {
		return core.SubmitJobResult{Job: *existing, Created: false}, nil
	}

	created, err := s.repo.CreateQueued(ctx, userID, req)
	if err != nil {
		if errors.Is(err, core.ErrConflict) {
			resolved, lookupErr := s.repo.FindByIdempotencyKey(ctx, userID, idempotencyKey)
			if lookupErr == nil && resolved != nil {
				return core.SubmitJobResult{Job: *resolved, Created: false}, nil
			}
		}
		return core.SubmitJobResult{}, err
	}

	enqueueErr := s.dispatcher.Enqueue(ctx, core.JobDispatchRequest{
		JobID:           created.Summary.JobID,
		UserID:          userID,
		SimulateFailure: created.SimulateFailure,
		IdempotencyKey:  idempotencyKey,
		CurrentAttempt:  created.Attempt,
	})
	if enqueueErr != nil {
		_ = s.repo.SetFailed(ctx, created.Summary.JobID, "failed to enqueue job")
		return core.SubmitJobResult{}, enqueueErr
	}

	resolved, err := s.repo.GetByID(ctx, userID, created.Summary.JobID)
	if err != nil {
		return core.SubmitJobResult{}, err
	}
	return core.SubmitJobResult{Job: *resolved, Created: true}, nil
}

func (s *JobService) ListJobs(ctx context.Context, userID string, limit, offset int) ([]core.JobSummary, error) {
	return s.repo.List(ctx, core.JobListFilter{UserID: userID, Limit: limit, Offset: offset})
}

func (s *JobService) GetJob(ctx context.Context, userID, jobID string) (*core.JobDetails, error) {
	if strings.TrimSpace(jobID) == "" {
		return nil, core.ErrInvalidArgument
	}
	return s.repo.GetByID(ctx, userID, jobID)
}

func (s *JobService) CancelJob(ctx context.Context, userID, jobID string) (*core.JobDetails, error) {
	job, err := s.repo.GetByID(ctx, userID, jobID)
	if err != nil {
		return nil, err
	}
	switch job.Summary.Status {
	case core.JobStatusQueued, core.JobStatusInProgress, core.JobStatusNew:
		return s.repo.SetCancelled(ctx, userID, jobID)
	default:
		return nil, core.ErrJobNotCancelable
	}
}

func (s *JobService) RetryJob(ctx context.Context, userID, jobID string) (*core.JobDetails, error) {
	job, err := s.repo.GetByID(ctx, userID, jobID)
	if err != nil {
		return nil, err
	}
	if job.Summary.Status != core.JobStatusFailed && job.Summary.Status != core.JobStatusCancelled {
		return nil, core.ErrJobNotRetryable
	}

	updated, err := s.repo.ResetForRetry(ctx, userID, jobID)
	if err != nil {
		return nil, err
	}
	enqueueErr := s.dispatcher.Enqueue(ctx, core.JobDispatchRequest{
		JobID:           updated.Summary.JobID,
		UserID:          userID,
		SimulateFailure: updated.SimulateFailure,
		IdempotencyKey:  derefString(updated.Summary.IdempotencyKey),
		CurrentAttempt:  updated.Attempt,
	})
	if enqueueErr != nil {
		_ = s.repo.SetFailed(ctx, updated.Summary.JobID, "failed to enqueue retry")
		return nil, enqueueErr
	}
	return s.repo.GetByID(ctx, userID, jobID)
}

func (s *JobService) GetJobResultURLs(
	ctx context.Context,
	userID, jobID string,
	ttl time.Duration,
) ([]core.ResultFileURL, error) {
	details, err := s.repo.GetByID(ctx, userID, jobID)
	if err != nil {
		return nil, err
	}
	urls := make([]core.ResultFileURL, 0, len(details.OutputFiles))
	for _, output := range details.OutputFiles {
		resolved, err := s.resolver.ResolveResultURL(ctx, output.Key, ttl)
		if err != nil {
			return nil, err
		}
		resolved.FileName = output.FileName
		urls = append(urls, resolved)
	}
	return urls, nil
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
