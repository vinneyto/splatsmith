package standalone

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type SimulatedReconstructionSubmissionDispatcher struct {
	repo core.ReconstructionJobRepository

	queue chan core.ReconstructionSubmissionRequest
	stop  chan struct{}
	wg    sync.WaitGroup
}

func NewSimulatedReconstructionSubmissionDispatcher(repo core.ReconstructionJobRepository) *SimulatedReconstructionSubmissionDispatcher {
	d := &SimulatedReconstructionSubmissionDispatcher{
		repo:  repo,
		queue: make(chan core.ReconstructionSubmissionRequest, 128),
		stop:  make(chan struct{}),
	}
	d.wg.Add(1)
	go d.loop()
	return d
}

func (d *SimulatedReconstructionSubmissionDispatcher) Enqueue(_ context.Context, req core.ReconstructionSubmissionRequest) error {
	select {
	case d.queue <- req:
		return nil
	default:
		return fmt.Errorf("simulated dispatcher queue is full: %w", core.ErrConflict)
	}
}

func (d *SimulatedReconstructionSubmissionDispatcher) Close() error {
	close(d.stop)
	d.wg.Wait()
	return nil
}

func (d *SimulatedReconstructionSubmissionDispatcher) loop() {
	defer d.wg.Done()
	for {
		select {
		case <-d.stop:
			return
		case req := <-d.queue:
			d.runSimulation(req)
		}
	}
}

func (d *SimulatedReconstructionSubmissionDispatcher) runSimulation(req core.ReconstructionSubmissionRequest) {
	ctx := context.Background()
	if err := d.repo.SetRunning(ctx, req.JobID); err != nil {
		return
	}

	steps := []struct {
		progress int
		step     string
		delay    time.Duration
	}{
		{progress: 10, step: "ingest", delay: 900 * time.Millisecond},
		{progress: 35, step: "preprocess", delay: 900 * time.Millisecond},
		{progress: 62, step: "reconstruct", delay: 1400 * time.Millisecond},
		{progress: 88, step: "postprocess", delay: 900 * time.Millisecond},
	}

	for i, st := range steps {
		time.Sleep(st.delay)
		if err := d.repo.SetProgress(ctx, req.JobID, st.progress, st.step); err != nil {
			return
		}
		if req.SimulateFailure && i >= 2 {
			_ = d.repo.SetFailed(ctx, req.JobID, "simulated failure in reconstruct step")
			return
		}
	}

	outputs := []core.OutputFileRef{
		{Key: fmt.Sprintf("outputs/%s/model.splat", req.JobID), FileName: "model.splat"},
		{Key: fmt.Sprintf("outputs/%s/model.ply", req.JobID), FileName: "model.ply"},
	}
	_ = d.repo.SetDone(ctx, req.JobID, outputs)
}
