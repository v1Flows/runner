package flow

import (
	"alertflow-runner/models"
)

func CheckFlowActions(flow models.Flows) (status bool) {
	// check if flow got any action
	if len(flow.Actions) == 0 {
		return false
	}

	return true
}
