package aws

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vinneyto/splatmaker/api/internal/core"
)

type S3ResultURLResolver struct {
	presign    *s3.PresignClient
	bucketName string
}

func NewS3ResultURLResolver(cfg Config) (*S3ResultURLResolver, error) {
	if strings.TrimSpace(cfg.ResultBucket) == "" {
		return nil, fmt.Errorf("aws.result_bucket is required")
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(awsCfg)
	return &S3ResultURLResolver{presign: s3.NewPresignClient(client), bucketName: cfg.ResultBucket}, nil
}

func (r *S3ResultURLResolver) ResolveResultURL(ctx context.Context, key string, ttl time.Duration) (core.ResultFileURL, error) {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	if strings.TrimSpace(key) == "" {
		return core.ResultFileURL{}, core.ErrInvalidArgument
	}
	res, err := r.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &r.bucketName,
		Key:    &key,
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return core.ResultFileURL{}, err
	}
	return core.ResultFileURL{
		Key:       key,
		FileName:  path.Base(key),
		URL:       res.URL,
		ExpiresAt: time.Now().UTC().Add(ttl),
	}, nil
}
