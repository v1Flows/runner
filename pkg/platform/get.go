package platform

import "sync"

// Global map to store platform information for each execution
var ExecutionPlatformMap = make(map[string]string)
var mu sync.Mutex

// Function to retrieve platform information for a given execution ID
func GetPlatformForExecution(executionID string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()
	platform, ok := ExecutionPlatformMap[executionID]
	return platform, ok
}

func SetPlatformForExecution(executionID, platform string) {
	mu.Lock()
	defer mu.Unlock()
	ExecutionPlatformMap[executionID] = platform
}
