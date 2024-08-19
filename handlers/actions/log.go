package actions

import "alertflow-runner/models"

func LogInit() models.ActionDetails {
	return models.ActionDetails{
		Name:        "Log Message",
		Description: "Prints an Log Message on the API Backend",
		Type:        "log",
		Params:      nil,
	}
}
