package standalone

import (
	"context"
	"fmt"
	"time"
)

type PipelineStub struct{}

func NewPipelineStub() *PipelineStub { return &PipelineStub{} }

func (p *PipelineStub) StartScan(_ context.Context, scanID, _ string) (string, error) {
	return fmt.Sprintf("standalone-job-%s-%d", scanID, time.Now().Unix()), nil
}
