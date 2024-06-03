package plugin

import (
	"alertflow-runner/src/config"
	"fmt"
	"log"
)

type RegisterActionArgs struct {
	Name string
	Type string
}

type Action string

func (a *Action) RegisterAction(args *RegisterActionArgs, reply *string) error {
	fmt.Println("Register Action: ", args.Name)

	config, err := config.ReadConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	var ApiURL = config.Alertflow.URL
	var ApiKey = config.Alertflow.APIKey
	var RunnerID = config.RunnerID

	fmt.Println(ApiURL, ApiKey, RunnerID)

	*reply = "Action Registered"
	return nil
}
