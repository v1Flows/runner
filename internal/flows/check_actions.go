package flows

import (
	"alertflow-runner/internal/executions"
	"alertflow-runner/pkg/models"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func CheckFlowActions(flow models.Flows, execution models.Execution) (status bool, err error) {
	stepData, err := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Check for Actions",
		ActionMessages: []string{"Checking if Flow has any Actions"},
		StartedAt:      time.Now(),
		Icon:           "solar:minimalistic-magnifer-linear",
	})
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	// check if flow got any action
	if len(flow.Actions) > 0 {
		count := 0
		for _, action := range flow.Actions {
			if action.Active {
				count++
			}
		}

		if count == 0 {
			err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
				ID:             stepData.ID,
				ActionMessages: []string{"Flow has no active Actions defined"},
				NoResult:       true,
				Finished:       true,
				FinishedAt:     time.Now(),
			})
			if err != nil {
				log.Error(fmt.Sprintf("Error updating step: %s", err))
				executions.EndWithError(execution)
				return false, err
			}
			executions.EndWithGhost(execution)
			return false, nil
		} else {
			err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
				ID:             stepData.ID,
				ActionMessages: []string{"Flow has Actions defined"},
				Finished:       true,
				FinishedAt:     time.Now(),
			})
			if err != nil {
				log.Error(fmt.Sprintf("Error updating step: %s", err))
				executions.EndWithError(execution)
				return false, err
			}
			return true, nil
		}
	} else {
		err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
			ID:             stepData.ID,
			ActionMessages: []string{"Flow has no Actions defined"},
			NoResult:       true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			log.Error(fmt.Sprintf("Error updating step: %s", err))
			executions.EndWithError(execution)
			return false, err
		}
		executions.EndWithGhost(execution)
		return false, nil
	}
}
