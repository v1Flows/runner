package actions

import (
	"alertflow-runner/functions/executions"
	"alertflow-runner/models"
	"time"
)

func StartAction(execution models.Execution, action models.FlowActions, payload models.Payload) (success bool, no_match bool, err error) {
	actionStepData, err := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Action: " + action.Name,
		ActionMessages: []string{"Starting Action: " + action.Name},
		StartedAt:      time.Now(),
	})
	if err != nil {
		return false, false, err
	}

	match, err := checkMatch(execution, action, payload, actionStepData.ID.String())
	if err != nil {
		return false, false, err
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
			return false, false, err
		}

		return false, true, nil
	}

	exclude, err := checkExclude(execution, action, payload, actionStepData.ID.String())
	if err != nil {
		return false, false, err
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
			return false, false, err
		}

		return false, true, nil
	}

	err = executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             actionStepData.ID,
		ActionMessages: []string{"Finished Action: " + action.Name},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		return false, false, err
	}

	return true, false, nil
}
