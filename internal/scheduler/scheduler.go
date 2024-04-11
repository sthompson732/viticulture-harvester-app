/*
 * File: scheduler.go
 * Description: Manages the scheduling of automated tasks for the Viniculture Data Harvester application.
 *              This module uses Google Cloud Scheduler to orchestrate tasks like
 *              satellite data retrieval, weather data updates, and soil data synchronization on a regular
 *              schedule.
 * Usage:
 *   - Configures and initiates scheduled jobs that trigger data retrieval and processing tasks.
 *   - Ensures tasks are executed at specified intervals, handling retries and logging as necessary.
 *   - Utilizes cron syntax to define job schedules.
 * Dependencies:
 *   - Requires external scheduling APIs or local cron services.
 *   - Interacts with client modules (e.g., satellite.go, weather.go, soil.go) to set up data fetch operations.
 *   - Uses service modules (e.g., imageservice.go, soildataservice.go) to process and store the fetched data.
 * Example:
 *   - A scheduler job might be set to retrieve satellite imagery every 12 hours and update soil data every 24 hours.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package scheduler

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/sthompson732/viticulture-harvester-app/internal/config"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	scheduler "cloud.google.com/go/scheduler/apiv1"
	schedulerpb "cloud.google.com/go/scheduler/apiv1/schedulerpb"
)

type SchedulerClient struct {
	Client *scheduler.CloudSchedulerClient
	Cfg    *config.Config
}

func NewSchedulerClient(ctx context.Context, cfg *config.Config) (*SchedulerClient, error) {
	client, err := scheduler.NewCloudSchedulerClient(ctx, option.WithCredentialsFile(cfg.CloudStorage.CredentialsPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler client: %w", err)
	}
	return &SchedulerClient{Client: client, Cfg: cfg}, nil
}

func (sc *SchedulerClient) InitializeJobs(ctx context.Context) error {
	for name, ds := range sc.Cfg.DataSources {
		if ds.Enabled {
			if err := sc.ensureJob(ctx, name, ds); err != nil {
				log.Printf("Failed to ensure job %s: %v", name, err)
			}
		}
	}
	return nil
}

func (sc *SchedulerClient) ensureJob(ctx context.Context, name string, ds config.DataSourceConfig) error {
	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", sc.Cfg.ProjectID, sc.Cfg.LocationID, name)
	_, err := sc.Client.GetJob(ctx, &schedulerpb.GetJobRequest{Name: jobName})

	if err != nil && status.Code(err) != codes.NotFound {
		return fmt.Errorf("error checking job %s: %w", name, err)
	}

	if err == nil {
		log.Printf("Job %s already exists", name)
		return nil
	}

	// Retry logic with exponential backoff
	retryCount := 3
	for i := 0; i < retryCount; i++ {
		if err = sc.createJob(ctx, jobName, ds); err == nil {
			log.Printf("Successfully created job %s on attempt %d", name, i+1)
			return nil
		}
		waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
		log.Printf("Failed to create job %s, retrying in %v...", name, waitTime)
		time.Sleep(waitTime)
	}

	return fmt.Errorf("failed to create job %s after %d retries: %v", name, retryCount, err)
}

func (sc *SchedulerClient) createJob(ctx context.Context, jobName string, ds config.DataSourceConfig) error {
	tz := ds.TimeZone
	if tz == "" {
		tz = "UTC" // Default to UTC if timezone is not specified
	}

	httpTarget := &schedulerpb.HttpTarget{
		Uri:        ds.Endpoint,
		HttpMethod: schedulerpb.HttpMethod(schedulerpb.HttpMethod_value[ds.HttpMethod]),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}

	jobReq := &schedulerpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", sc.Cfg.ProjectID, sc.Cfg.LocationID),
		Job: &schedulerpb.Job{
			Name:        jobName,
			Description: ds.Description,
			Schedule:    ds.Schedule,
			TimeZone:    tz,
			Target: &schedulerpb.Job_HttpTarget{ // This assumes 'Target' is the correct oneof field name for different job types
				HttpTarget: httpTarget,
			},
		},
	}

	_, err := sc.Client.CreateJob(ctx, jobReq)
	if err != nil {
		return fmt.Errorf("failed to create job %s: %v", jobName, err)
	}

	return nil
}
