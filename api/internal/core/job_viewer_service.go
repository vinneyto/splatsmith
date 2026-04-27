package core

import (
	"context"
	"time"
)

type JobViewerService struct {
	repo     JobRepository
	resolver ResultURLResolver
}

func NewJobViewerService(repo JobRepository, resolver ResultURLResolver) *JobViewerService {
	return &JobViewerService{repo: repo, resolver: resolver}
}

func (s *JobViewerService) ListJobs(ctx context.Context, userID string, limit, offset int) ([]JobSummary, error) {
	return s.repo.List(ctx, JobListFilter{UserID: userID, Limit: limit, Offset: offset})
}

func (s *JobViewerService) GetJobResultURLs(
	ctx context.Context,
	userID, jobID string,
	ttl time.Duration,
) ([]ResultFileURL, error) {
	details, err := s.repo.GetByID(ctx, userID, jobID)
	if err != nil {
		return nil, err
	}
	urls := make([]ResultFileURL, 0, len(details.OutputFiles))
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
