package actions

import (
	"encoding/json"
	"time"

	"gitlab.justlab.xyz/alertflow-public/runner/internal/executions"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"

	"github.com/tidwall/gjson"

	log "github.com/sirupsen/logrus"
)

func PatternCheckInit() models.ActionDetails {
	return models.ActionDetails{
		Name:        "Pattern Check",
		Description: "Check flow patterns",
		Icon:        "solar:list-check-minimalistic-bold",
		Type:        "pattern_check",
		Category:    "Utility",
		Function:    PatternCheckAction,
		IsHidden:    true,
		Params:      nil,
	}
}

func PatternCheckAction(execution models.Execution, flow models.Flows, payload models.Payload, allSteps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool) {
	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionID:       action.ID.String(),
		ActionMessages: []string{"Checking for patterns"},
		Pending:        false,
		Running:        true,
		StartedAt:      time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
	}

	// end if there are no patterns
	if len(flow.Patterns) == 0 {
		err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
			ID:             step.ID,
			ActionMessages: []string{"No patterns are defined. Continue to next step"},
			Running:        false,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			return nil, false, false, false, true
		}
		return nil, true, false, false, false
	}

	// convert payload to string
	payloadBytes, err := json.Marshal(payload.Payload)
	if err != nil {
		log.Error("Error converting payload to JSON:", err)
		return nil, false, false, false, true
	}
	payloadString := string(payloadBytes)

	patternMissMatched := 0

	for _, pattern := range flow.Patterns {
		value := gjson.Get(payloadString, pattern.Key)

		if pattern.Type == "equals" {
			if value.String() == pattern.Value {
				err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
					ID:             step.ID,
					ActionMessages: []string{`Pattern: ` + pattern.Key + ` == ` + pattern.Value + ` matched. Continue to next step`},
				})
				if err != nil {
					return nil, false, false, false, true
				}
			} else {
				err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
					ID:             step.ID,
					ActionMessages: []string{`Pattern: ` + pattern.Key + ` == ` + pattern.Value + ` not found.`},
					Running:        false,
					Canceled:       true,
					Finished:       true,
					FinishedAt:     time.Now(),
				})
				if err != nil {
					return nil, false, false, false, true
				}
				patternMissMatched++
			}
		} else if pattern.Type == "not_equals" {
			if value.String() != pattern.Value {
				err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
					ID:             step.ID,
					ActionMessages: []string{`Pattern: ` + pattern.Key + ` != ` + pattern.Value + ` not found. Continue to next step`},
				})
				if err != nil {
					return nil, false, false, false, true
				}
			} else {
				err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
					ID:             step.ID,
					ActionMessages: []string{`Pattern: ` + pattern.Key + ` != ` + pattern.Value + ` matched.`},
					Running:        false,
					Canceled:       true,
					Finished:       true,
					FinishedAt:     time.Now(),
				})
				if err != nil {
					return nil, false, false, false, true
				}
				patternMissMatched++
			}
		}
	}

	if patternMissMatched > 0 {
		err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
			ID:             step.ID,
			ActionMessages: []string{"Some patterns did not match. Cancel execution"},
			Running:        false,
			Canceled:       false,
			NoPatternMatch: true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			return nil, false, false, false, true
		}
		return nil, false, false, true, false
	} else {
		err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
			ID:             step.ID,
			ActionMessages: []string{"All patterns matched. Continue to next step"},
			Running:        false,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			return nil, false, false, false, true
		}
		return nil, true, false, false, false
	}
}
