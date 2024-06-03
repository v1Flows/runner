package plugin

import (
	"fmt"
)

type RegisterActionArgs struct {
	Name string
	Type string
}

type Action string

func (a *Action) RegisterAction(args *RegisterActionArgs, reply *string) error {
	fmt.Println("Register Action: ", args)
	*reply = "Action Registered"
	return nil
}
