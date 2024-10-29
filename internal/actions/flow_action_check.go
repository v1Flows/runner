package actions

import (
	"alertflow-runner/internal/executions"
	"alertflow-runner/pkg/models"
	"time"
)

func FlowActionsCheckInit() models.ActionDetails {
	return models.ActionDetails{
		Name:        "Flow Action Check",
		Description: "Check flow for actions",
		Icon:        "solar:minimalistic-magnifer-linear",
		Type:        "flow_actions_check",
		Function:    FlowActionsCheckAction,
		IsHidden:    true,
		Params:      nil,
	}
}

func FlowActionsCheckAction(execution models.Execution, flow models.Flows, payload models.Payload, allSteps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool) {
	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionID:       action.ID.String(),
		ActionMessages: []string{"Checking for flow actions"},
		Pending:        false,
		Running:        true,
		StartedAt:      time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
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
				ID:             step.ID,
				ActionMessages: []string{"Flow has no active Actions defined. Cancel execution"},
				Running:        false,
				Canceled:       true,
				CanceledBy:     "Flow Action Check",
				CanceledAt:     time.Now(),
				Finished:       true,
				FinishedAt:     time.Now(),
			})
			if err != nil {
				return nil, false, false, false, true
			}
			return nil, false, true, false, false
		} else {
			err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
				ID:             step.ID,
				ActionMessages: []string{"Flow has Actions defined"},
				Running:        false,
				Finished:       true,
				FinishedAt:     time.Now(),
			})
			if err != nil {
				return nil, false, false, false, true
			}
			return nil, true, false, false, false
		}
	} else {
		err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
			ID:             step.ID,
			ActionMessages: []string{"Flow has no Actions defined. Cancel execution"},
			Canceled:       true,
			CanceledBy:     "Flow Action Check",
			CanceledAt:     time.Now(),
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			return nil, false, false, false, true
		}
		return nil, false, true, false, false
	}
}
