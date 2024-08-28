package actions

import (
	"alertflow-runner/functions/executions"
	handler_actions "alertflow-runner/handlers/actions"
	"alertflow-runner/models"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func StartAction(execution models.Execution, action models.FlowActions, payload models.Payload) (success bool, no_match bool, action_error bool, err error) {
	actionStepData, err := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Action: " + action.Name,
		ActionMessages: []string{"Starting Action: " + action.Name},
		StartedAt:      time.Now(),
		Icon:           "solar:bolt-line-duotone",
	})
	if err != nil {
		return false, false, false, err
	}

	match, err := checkMatch(execution, action, payload, actionStepData.ID.String())
	if err != nil {
		return false, false, false, err
	}

	if !match {
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             actionStepData.ID,
			ActionMessages: []string{"Match pattern not found for: " + action.Name},
			NoPatternMatch: true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			return false, false, false, err
		}

		return false, true, false, nil
	}

	exclude, err := checkExclude(execution, action, payload, actionStepData.ID.String())
	if err != nil {
		return false, false, false, err
	}

	if exclude {
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             actionStepData.ID,
			ActionMessages: []string{"Exclude pattern found for: " + action.Name},
			NoPatternMatch: true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			return false, false, false, err
		}

		return false, true, false, nil
	}

	var actionErrors []string

	// check for parallel or sequential order
	if action.ExecParallel {
		fmt.Println("Action: " + action.Name + " is set to run in parallel")
	} else {
		for _, subAction := range action.Actions {
			// search for the action
			subActionStepData, err := executions.SendStep(execution, models.ExecutionSteps{
				ExecutionID:    execution.ID.String(),
				ActionName:     subAction,
				ActionMessages: []string{"Starting: " + subAction},
				StartedAt:      time.Now(),
				ParentID:       actionStepData.ID.String(),
				IsHidden:       true,
			})
			if err != nil {
				return false, false, false, err
			}
			actionDetails := handler_actions.SearchAction(subAction)

			if actionDetails.Name == "" {
				err = executions.UpdateStep(execution, models.ExecutionSteps{
					ID:             subActionStepData.ID,
					ActionMessages: []string{"Action: " + subAction + " not found"},
					Error:          true,
					Finished:       true,
					FinishedAt:     time.Now(),
				})
				if err != nil {
					return false, false, true, err
				}
				continue
			} else {
				// exec the actionDetails.Function
				if fn, ok := actionDetails.Function.(func()); ok {
					fn()
				} else {
					// handle the case when actionDetails.Function is not a function
					err = executions.UpdateStep(execution, models.ExecutionSteps{
						ID:             subActionStepData.ID,
						ActionMessages: []string{"Action: " + subAction + " is not a function"},
						Error:          true,
						Finished:       true,
						FinishedAt:     time.Now(),
					})
					if err != nil {
						log.Error(err)
					}
					actionErrors = append(actionErrors, subAction)
					break
				}
			}
		}
	}

	if len(actionErrors) > 0 {
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             actionStepData.ID,
			ActionMessages: []string{"Error: " + fmt.Sprintf("%v", actionErrors)},
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			return false, false, false, err
		}
		return false, false, true, nil
	} else {
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             actionStepData.ID,
			ActionMessages: []string{"Action: " + action.Name + " completed successfully"},
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			return false, false, false, err
		}
		return true, false, false, nil
	}
}
