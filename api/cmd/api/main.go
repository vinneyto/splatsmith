package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/vinneyto/ariadne/api/internal/app"
)

func main() {
	configPath := flag.String("config", "./config/standalone.yaml", "path to YAML config")
	token := flag.String("token", "dev-token", "token for startup auth probe")
	flag.Parse()

	cfg, err := app.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	runtime, err := app.BuildRuntime(cfg)
	if err != nil {
		log.Fatalf("build runtime: %v", err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			log.Printf("runtime close error: %v", err)
		}
	}()

	identity, err := runtime.AuthService.Authenticate(context.Background(), "Bearer "+*token)
	if err != nil {
		log.Printf("startup auth probe failed (expected for aws stubs): %v", err)
		fmt.Printf("api bootstrap is ready (mode=%s)\n", runtime.Mode)
		return
	}

	jobs, err := runtime.JobViewer.ListJobs(context.Background(), identity.UserID, 20, 0)
	if err != nil {
		log.Printf("job list failed (expected for aws stubs): %v", err)
		fmt.Printf("api bootstrap is ready (mode=%s)\n", runtime.Mode)
		return
	}

	fmt.Printf("api bootstrap is ready (mode=%s, user=%s, jobs=%d)\n", runtime.Mode, identity.UserID, len(jobs))
	for _, job := range jobs {
		urls, err := runtime.JobViewer.GetJobResultURLs(context.Background(), identity.UserID, job.JobID, time.Duration(runtime.ResultURLTTL)*time.Second)
		if err != nil {
			log.Printf("resolve urls for job %s failed: %v", job.JobID, err)
			continue
		}
		fmt.Printf("- job=%s status=%s files=%d\n", job.JobID, job.Status, len(urls))
		for _, u := range urls {
			fmt.Printf("  * %s -> %s\n", u.FileName, u.URL)
		}
	}
}
