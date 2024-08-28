package handler_actions

import "alertflow-runner/models"

func SearchAction(actionName string) models.ActionDetails {
	// Search for action
	// If action found, return action details
	// If action not found, return nil

	// get all actions
	actions := Init()

	for _, action := range actions {
		if action.Name == actionName {
			return action
		}
	}

	return models.ActionDetails{}
}
