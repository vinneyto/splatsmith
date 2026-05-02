package services

import (
	"context"
	"strings"
	"time"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type JobService struct {
	repo     core.JobRepository
	resolver core.ResultURLResolver
}

func NewJobService(repo core.JobRepository, resolver core.ResultURLResolver) *JobService {
	return &JobService{repo: repo, resolver: resolver}
}

func (s *JobService) ListJobs(ctx context.Context, filter core.JobListFilter) ([]core.JobSummary, error) {
	return s.repo.List(ctx, filter)
}

func (s *JobService) GetJob(ctx context.Context, jobID string) (*core.JobDetails, error) {
	if strings.TrimSpace(jobID) == "" {
		return nil, core.ErrInvalidArgument
	}
	return s.repo.GetByID(ctx, jobID)
}

func (s *JobService) GetJobResultURLs(ctx context.Context, jobID string, ttl time.Duration) ([]core.ResultFileURL, error) {
	details, err := s.repo.GetByID(ctx, jobID)
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
