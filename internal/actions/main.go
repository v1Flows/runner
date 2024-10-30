package actions

import (
	"alertflow-runner/internal/executions"
	"alertflow-runner/pkg/models"
	"time"

	log "github.com/sirupsen/logrus"
)

func StartAction(action models.Actions, execution models.Execution, steps []models.ExecutionSteps) (data any, finished bool, canceled bool, failed bool, err error) {
	_, err = executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     action.Name,
		ActionMessages: []string{"Starting Action: " + action.Name + " | ID: " + action.ID.String()},
		StartedAt:      time.Now(),
		Icon:           action.Icon,
	})
	if err != nil {
		log.Error(err)
		return nil, false, false, false, err
	}

	// if action.Name == "" {
	// 	err = executions.UpdateStep(execution, models.ExecutionSteps{
	// 		ID:             actionStepData.ID,
	// 		ActionMessages: []string{"Action: " + action.Name + " not found"},
	// 		Error:          true,
	// 		Finished:       true,
	// 		FinishedAt:     time.Now(),
	// 	})
	// 	if err != nil {
	// 		log.Error(err)
	// 		return nil, false, false, false, err
	// 	}
	// 	return nil, false, false, true, nil
	// } else {
	// 	// exec the actionDetails.Function
	// 	if fn, ok := action.Function.(func(execution models.Execution, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data any, finished bool, canceled bool, failed bool)); ok {
	// 		_, finished, canceled, failed := fn(execution, steps, actionStepData, action)

	// 		if failed {
	// 			err = executions.UpdateStep(execution, models.ExecutionSteps{
	// 				ID:             actionStepData.ID,
	// 				ActionMessages: []string{"Action: " + action.Name + " failed"},
	// 				Error:          true,
	// 				Finished:       true,
	// 				FinishedAt:     time.Now(),
	// 			})
	// 			if err != nil {
	// 				log.Error(err)
	// 				return nil, false, false, false, err
	// 			}
	// 			return nil, false, false, true, nil
	// 		} else if canceled {
	// 			err = executions.UpdateStep(execution, models.ExecutionSteps{
	// 				ID:             actionStepData.ID,
	// 				ActionMessages: []string{"Action: " + action.Name + " canceled"},
	// 				Error:          true,
	// 				Finished:       true,
	// 				FinishedAt:     time.Now(),
	// 			})
	// 			if err != nil {
	// 				log.Error(err)
	// 				return nil, false, false, false, err
	// 			}
	// 			return nil, false, true, false, nil
	// 		} else if finished {
	// 			return nil, true, false, false, nil
	// 		}
	// 	} else {
	// 		// handle the case when actionDetails.Function is not a function
	// 		err = executions.UpdateStep(execution, models.ExecutionSteps{
	// 			ID:             actionStepData.ID,
	// 			ActionMessages: []string{"Action: " + action.Name + " is not a function"},
	// 			Error:          true,
	// 			Finished:       true,
	// 			FinishedAt:     time.Now(),
	// 		})
	// 		if err != nil {
	// 			log.Error(err)
	// 			return nil, false, false, false, err
	// 		}
	// 		return nil, false, false, true, nil
	// 	}
	// }

	return nil, true, false, false, nil
}
