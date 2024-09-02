package actions

import (
	"alertflow-runner/functions/executions"
	handler_actions "alertflow-runner/handlers/actions"
	"alertflow-runner/models"
	"time"

	log "github.com/sirupsen/logrus"
)

func StartAction(action models.Actions, execution models.Execution) (status bool, err error) {
	actionStepData, err := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     action.Name,
		ActionMessages: []string{"Starting Action: " + action.Name + " | ID: " + action.ID.String()},
		StartedAt:      time.Now(),
		Icon:           action.Icon,
	})
	if err != nil {
		log.Error(err)
		return false, err
	}

	actionDetails := handler_actions.SearchAction(action.Name)

	if actionDetails.Name == "" {
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             actionStepData.ID,
			ActionMessages: []string{"Action: " + action.Name + " not found"},
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			log.Error(err)
			return false, err
		}
		return false, nil
	} else {
		// exec the actionDetails.Function
		if fn, ok := actionDetails.Function.(func(step models.ExecutionSteps, action models.Actions) bool); ok {
			status := fn(actionStepData, action)

			if !status {
				err = executions.UpdateStep(execution, models.ExecutionSteps{
					ID:             actionStepData.ID,
					ActionMessages: []string{"Action: " + action.Name + " failed"},
					Error:          true,
					Finished:       true,
					FinishedAt:     time.Now(),
				})
				if err != nil {
					log.Error(err)
					return false, err
				}
				return false, nil
			}
		} else {
			// handle the case when actionDetails.Function is not a function
			err = executions.UpdateStep(execution, models.ExecutionSteps{
				ID:             actionStepData.ID,
				ActionMessages: []string{"Action: " + action.Name + " is not a function"},
				Error:          true,
				Finished:       true,
				FinishedAt:     time.Now(),
			})
			if err != nil {
				log.Error(err)
				return false, err
			}
			return false, nil
		}
	}

	return true, nil
}
