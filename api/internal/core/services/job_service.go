package services

import (
	"context"
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

func (s *JobService) ListJobs(ctx context.Context, userID string, limit, offset int) ([]core.JobSummary, error) {
	return s.repo.List(ctx, core.JobListFilter{UserID: userID, Limit: limit, Offset: offset})
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
