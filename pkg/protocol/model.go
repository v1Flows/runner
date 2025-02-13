package protocol

type Request struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

type Response struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Error   string                 `json:"error,omitempty"`
}
