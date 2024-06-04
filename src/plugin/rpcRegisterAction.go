package plugin

import (
	"bytes"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type RegisterActionArgs struct {
	Name        string
	Type        string
	Description string
	Fields      json.RawMessage
}

type Action string

func (a *Action) RegisterAction(args *RegisterActionArgs, reply *string) error {
	reader := bytes.NewReader(args.Fields)
	decoder := json.NewDecoder(reader)

	var fields map[string]interface{}
	err := decoder.Decode(&fields)
	if err != nil {
		return err
	}

	log.Info("Action Register incoming: Name ", args.Name, "; Description ", args.Description, "; Type ", args.Type, "; Fields ", fields)

	*reply = "Action Registered"
	return nil
}
