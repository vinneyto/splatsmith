package aws

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/vinneyto/splatmaker/api/internal/core"
)

type DynamoJobRepository struct {
	ddb       *dynamodb.Client
	tableName string
}

type dynamoJobItem struct {
	UUID         string `dynamodbav:"uuid"`
	UUIDStatus   string `dynamodbav:"uuidStatus"`
	StartTS      string `dynamodbav:"startTimestamp"`
	EndTS        string `dynamodbav:"endTimestamp"`
	UpdateTS     string `dynamodbav:"updatedAt"`
	S3Output     string `dynamodbav:"s3Output"`
	S3Input      string `dynamodbav:"s3Input"`
	Filename     string `dynamodbav:"filename"`
	ErrorMessage string `dynamodbav:"errorMsg"`
}

func NewDynamoJobRepository(cfg Config) (*DynamoJobRepository, error) {
	if strings.TrimSpace(cfg.JobsTable) == "" {
		return nil, fmt.Errorf("aws.jobs_table is required")
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, err
	}
	return &DynamoJobRepository{ddb: dynamodb.NewFromConfig(awsCfg), tableName: cfg.JobsTable}, nil
}

func (r *DynamoJobRepository) List(ctx context.Context, filter core.JobListFilter) ([]core.JobSummary, error) {
	out, err := r.ddb.Scan(ctx, &dynamodb.ScanInput{TableName: aws.String(r.tableName)})
	if err != nil {
		return nil, err
	}
	items := make([]core.JobSummary, 0, len(out.Items))
	for _, raw := range out.Items {
		var row dynamoJobItem
		if err := attributevalue.UnmarshalMap(raw, &row); err != nil {
			continue
		}
		items = append(items, toCoreSummary(row, filter.UserID))
	}
	sort.Slice(items, func(i, j int) bool { return items[i].UpdatedAt.After(items[j].UpdatedAt) })
	if filter.Offset >= len(items) {
		return []core.JobSummary{}, nil
	}
	end := len(items)
	if filter.Limit > 0 && filter.Offset+filter.Limit < end {
		end = filter.Offset + filter.Limit
	}
	return items[filter.Offset:end], nil
}

func (r *DynamoJobRepository) GetByID(ctx context.Context, userID, jobID string) (*core.JobDetails, error) {
	out, err := r.ddb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       map[string]types.AttributeValue{"uuid": &types.AttributeValueMemberS{Value: jobID}},
	})
	if err != nil {
		return nil, err
	}
	if len(out.Item) == 0 {
		return nil, core.ErrJobNotFound
	}
	var row dynamoJobItem
	if err := attributevalue.UnmarshalMap(out.Item, &row); err != nil {
		return nil, err
	}
	summary := toCoreSummary(row, userID)
	outputs := inferOutputFiles(row)
	return &core.JobDetails{
		Summary:      summary,
		OutputFiles:  outputs,
		Attempt:      1,
		SourceRef:    stringPtr(strings.TrimSpace(row.S3Input)),
		ErrorMessage: stringPtr(strings.TrimSpace(row.ErrorMessage)),
		StartedAt:    parseTimePtr(row.StartTS),
		FinishedAt:   parseTimePtr(row.EndTS),
	}, nil
}

func toCoreSummary(in dynamoJobItem, userID string) core.JobSummary {
	status := mapStatus(in.UUIDStatus)
	updated := parseTime(in.UpdateTS)
	if updated.IsZero() {
		if x := parseTime(in.EndTS); !x.IsZero() {
			updated = x
		} else {
			updated = parseTime(in.StartTS)
		}
	}
	if updated.IsZero() {
		updated = time.Now().UTC()
	}
	created := parseTime(in.StartTS)
	if created.IsZero() {
		created = updated
	}
	progress := inferProgress(status)
	return core.JobSummary{
		JobID:           in.UUID,
		UserID:          userID,
		Status:          status,
		ProgressPercent: progress,
		CreatedAt:       created,
		UpdatedAt:       updated,
	}
}

func inferOutputFiles(in dynamoJobItem) []core.OutputFileRef {
	prefix := strings.TrimSpace(in.S3Output)
	if prefix == "" {
		return nil
	}
	clean := strings.TrimPrefix(prefix, "s3://")
	parts := strings.SplitN(clean, "/", 2)
	if len(parts) != 2 {
		return nil
	}
	base := strings.TrimSuffix(parts[1], "/")
	jobPart := strings.TrimSpace(in.UUID)
	if !strings.Contains(base, jobPart) {
		base = strings.TrimSuffix(base, "/") + "/" + jobPart
	}
	return []core.OutputFileRef{
		{Key: base + "/model.splat", FileName: "model.splat"},
		{Key: base + "/model.ply", FileName: "model.ply"},
		{Key: base + "/model.spz", FileName: "model.spz"},
	}
}

func mapStatus(v string) core.JobStatus {
	s := strings.ToLower(strings.TrimSpace(v))
	switch s {
	case "done", "completed", "success", "succeeded":
		return core.JobStatusDone
	case "failed", "error":
		return core.JobStatusFailed
	case "cancelled", "canceled":
		return core.JobStatusCancelled
	case "in-progress", "in_progress", "running", "processing":
		return core.JobStatusInProgress
	case "queued", "pending":
		return core.JobStatusQueued
	default:
		return core.JobStatusNew
	}
}

func inferProgress(status core.JobStatus) int {
	switch status {
	case core.JobStatusDone:
		return 100
	case core.JobStatusFailed, core.JobStatusCancelled:
		return 0
	case core.JobStatusInProgress:
		return 50
	case core.JobStatusQueued:
		return 10
	default:
		return 0
	}
}

func parseTime(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	formats := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05"}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t.UTC()
		}
	}
	if unix, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(unix, 0).UTC()
	}
	return time.Time{}
}

func parseTimePtr(s string) *time.Time {
	t := parseTime(s)
	if t.IsZero() {
		return nil
	}
	return &t
}

func stringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
