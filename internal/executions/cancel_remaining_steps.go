package internal_executions

import (
	"time"

	"github.com/AlertFlow/runner/pkg/executions"
	log "github.com/sirupsen/logrus"
)

func CancelRemainingSteps(executionID string) error {
	// get all steps where pending is true
	steps, err := executions.GetSteps(executionID)
	if err != nil {
		log.Error(err)
		return err
	}

	// cancel each step where pending is true
	for _, step := range steps {
		if step.Pending {
			step.Pending = false
			step.Canceled = true
			step.CanceledBy = "Runner"
			step.CanceledAt = time.Now()
			step.ActionMessages = []string{"Canceled by runner due to previous step failure/interaction/timeout"}
			step.StartedAt = time.Now()
			step.FinishedAt = time.Now()

			err := executions.UpdateStep(executionID, step)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}

	return nil
}
