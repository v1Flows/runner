package flow

import (
	"alertflow-runner/models"
)

func CheckFlowActions(actions []models.FlowActions) (status bool) {
	// check if flow got any action
	if len(actions) > 0 {
		return true
	}

	return false
}
