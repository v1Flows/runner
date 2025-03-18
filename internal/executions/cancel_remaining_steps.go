package internal_executions

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/executions"
	"github.com/v1Flows/runner/pkg/platform"
	shared_models "github.com/v1Flows/shared-library/pkg/models"
)

func cancelRemainingSteps(cfg config.Config, executionID string) error {
	targetPlatform, ok := platform.GetPlatformForExecution(executionID)
	if !ok {
		log.Error("Failed to get platform")
		return errors.New("failed to get platform")
	}

	steps, err := executions.GetSteps(cfg, executionID, targetPlatform)
	if err != nil {
		log.Error(err)
		return err
	}

	// cancel each step where pending is true
	for _, step := range steps {
		if step.Status == "pending" {
			step.Status = "canceled"
			step.CanceledBy = "Runner"
			step.CanceledAt = time.Now()
			step.Messages = []shared_models.Message{
				{
					Title: "Canceled",
					Lines: []string{"Canceled by runner due to previous step failure/interaction/timeout"},
				},
			}
			step.StartedAt = time.Now()
			step.FinishedAt = time.Now()

			err := executions.UpdateStep(cfg, executionID, step, targetPlatform)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}

	return nil
}
