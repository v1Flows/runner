package plugins

// Define the methods plugins can call on the runner
type RunnerRPC interface {
	NotifyStatus(status string, _ *struct{}) error
}

// Implementation of the runner-side RPC server
type RunnerRPCServer struct{}

func (r *RunnerRPCServer) NotifyStatus(status string, _ *struct{}) error {
	// Handle the callback from the plugin (e.g., log or update state)
	println("Plugin callback received:", status)
	return nil
}
