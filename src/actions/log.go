package actions

import "encoding/json"

type ActionDetails struct {
	Name        string
	Description string
	Type        string
	Params      json.RawMessage
}
