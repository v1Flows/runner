package executions

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func CancelRemainingSteps(executionID string) error {
	// get all steps where pending is true
	steps, err := GetSteps(executionID)
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
			step.ActionMessages = []string{"Canceled by runner due to previous step failure"}
			step.FinishedAt = time.Now()

			err := UpdateStep(executionID, step)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}

	return nil
}
