/*
 * scheduler.go: Orchestrates timed data fetching tasks using Google Cloud Scheduler.
 * Usage:
 *   - Configures and initiates scheduled jobs that trigger data retrieval and processing tasks.
 *   - Ensures tasks are executed at specified intervals, handling retries and logging as necessary.
 *   - Utilizes cron syntax to define job schedules.
* Dependencies:
 *   - Requires external scheduling APIs or local cron services.
 *   - Interacts with client modules (e.g., satellite.go, weather.go, soil.go) to set up data fetch operations.
 *   - Uses service modules (e.g., imageservice.go, soildataservice.go) to process and store the fetched data.
 * Author(s): Shannon Thompson
 * Created on: 04/12/2024
*/
package scheduler

import (
	"context"
	"fmt"
	"log"
	"strings"

	scheduler "cloud.google.com/go/scheduler/apiv1"
	"cloud.google.com/go/scheduler/apiv1/schedulerpb"
	"github.com/sthompson732/viticulture-harvester-app/internal/config"
	"google.golang.org/api/option"
)

type SchedulerClient struct {
	Client *scheduler.CloudSchedulerClient
	Cfg    *config.Config
}

func NewSchedulerClient(ctx context.Context, cfg *config.Config) (*SchedulerClient, error) {
	client, err := scheduler.NewCloudSchedulerClient(ctx, option.WithCredentialsFile(cfg.CloudStorage.CredentialsPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler client: %v", err)
	}
	return &SchedulerClient{
		Client: client,
		Cfg:    cfg,
	}, nil
}

func (sc *SchedulerClient) SetupJobs(ctx context.Context) error {
	for _, jobCfg := range sc.Cfg.DataSources {
		if jobCfg.Enabled {
			err := sc.createJob(ctx, jobCfg)
			if err != nil {
				log.Printf("Failed to create job for %s: %v", jobCfg.Description, err)
				continue
			}
			log.Printf("Successfully scheduled job: %s", jobCfg.Description)
		}
	}
	return nil
}

func (sc *SchedulerClient) createJob(ctx context.Context, jobCfg config.DataSourceConfig) error {
	parent := fmt.Sprintf("projects/%s/locations/%s", sc.Cfg.ProjectID, sc.Cfg.LocationID)

	// Build the HTTP target based on the new documentation
	httpTarget := &schedulerpb.HttpTarget{
		Uri:        jobCfg.Endpoint,
		HttpMethod: schedulerpb.HttpMethod(schedulerpb.HttpMethod_value[jobCfg.HttpMethod]),
	}

	// Add headers if any
	if len(jobCfg.Headers) > 0 {
		httpTarget.Headers = jobCfg.Headers
	}

	// Set the body if the method is POST, PUT, or PATCH
	if jobCfg.HttpMethod == "POST" || jobCfg.HttpMethod == "PUT" || jobCfg.HttpMethod == "PATCH" {
		httpTarget.Body = []byte(jobCfg.Body)
	}

	// OAuthToken and OidcToken should be set if needed here

	job := &schedulerpb.Job{
		Name:     fmt.Sprintf("%s/jobs/%s", parent, formatJobName(jobCfg.Description)),
		Target:   &schedulerpb.Job_HttpTarget{HttpTarget: httpTarget},
		Schedule: jobCfg.Schedule,
		TimeZone: jobCfg.TimeZone,
	}

	// Use the CreateJob method of the Cloud Scheduler client
	_, err := sc.Client.CreateJob(ctx, &schedulerpb.CreateJobRequest{
		Parent: parent,
		Job:    job,
	})
	if err != nil {
		return fmt.Errorf("failed to create job for %s: %v", jobCfg.Description, err)
	}
	return nil
}

func buildQueryParams(params map[string]string) string {
	var parts []string
	for key, value := range params {
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(parts, "&")
}

func formatJobName(description string) string {
	return strings.ReplaceAll(strings.ToLower(description), " ", "-")
}

func (sc *SchedulerClient) DeleteJob(ctx context.Context, jobName string) error {
	// The DeleteJob call returns only an error.
	err := sc.Client.DeleteJob(ctx, &schedulerpb.DeleteJobRequest{
		Name: jobName,
	})
	if err != nil {
		// Handle the error properly.
		return fmt.Errorf("failed to delete job %s: %v", jobName, err)
	}
	return nil
}
