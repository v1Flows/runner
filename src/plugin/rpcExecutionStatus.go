package plugin

type ExecutionStatus struct {
	ID        string
	Status    string
	Failed    bool
	Remaining int
}

type Response string

func (r *Response) ExecutionStatus(args *ExecutionStatus, reply *string) error {
	*reply = "ok"
	return nil
}
