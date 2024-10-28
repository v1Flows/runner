package actions

import "alertflow-runner/pkg/models"

func SearchAction(actionID string) (action models.ActionDetails, found bool) {
	// Search for action
	// If action found, return action details
	// If action not found, return nil

	// get all actions
	actions := Init()

	for _, action := range actions {
		if action.ID == actionID {
			return action, true
		}
	}

	return models.ActionDetails{}, false
}
