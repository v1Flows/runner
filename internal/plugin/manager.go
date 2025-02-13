package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/protocol"
)

type PluginProcess struct {
	cmd    *exec.Cmd
	stdin  *json.Encoder
	stdout *json.Decoder
	mutex  sync.Mutex
}

type Manager struct {
	pluginDir     string
	pluginTempDir string
	plugins       map[string]*PluginProcess
	mutex         sync.RWMutex
}

func NewManager(pluginDir, pluginTempDir string) *Manager {
	return &Manager{
		pluginDir:     pluginDir,
		pluginTempDir: pluginTempDir,
		plugins:       make(map[string]*PluginProcess),
	}
}

func (m *Manager) DownloadPlugin(plugin config.PluginConf) error {
	repoPath := filepath.Join(m.pluginTempDir, plugin.Name)

	// Clone and build as a standalone executable instead of a .so file
	cmd := exec.Command("git", "clone", plugin.Url, repoPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone plugin repository: %w", err)
	}

	if plugin.Version != "" {
		cmd = exec.Command("git", "-C", repoPath, "checkout", plugin.Version)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to checkout version: %w", err)
		}
	}

	// Build as standalone executable
	outputPath := filepath.Join(m.pluginDir, plugin.Name)
	cmd = exec.Command("go", "build", "-o", outputPath)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build plugin: %w", err)
	}

	return nil
}

func (m *Manager) StartPlugin(plugin config.PluginConf) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	pluginPath := filepath.Join(m.pluginDir, plugin.Name)
	cmd := exec.Command(pluginPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start plugin: %w", err)
	}

	process := &PluginProcess{
		cmd:    cmd,
		stdin:  json.NewEncoder(stdin),
		stdout: json.NewDecoder(stdout),
		mutex:  sync.Mutex{},
	}

	m.plugins[plugin.Name] = process
	return nil
}

func (m *Manager) ExecutePlugin(pluginName string, req protocol.Request) (*protocol.Response, error) {
	m.mutex.RLock()
	process, exists := m.plugins[pluginName]
	m.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginName)
	}

	process.mutex.Lock()
	defer process.mutex.Unlock()

	if err := process.stdin.Encode(req); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var resp protocol.Response
	if err := process.stdout.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return &resp, nil
}

func (m *Manager) StopPlugin(pluginName string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	process, exists := m.plugins[pluginName]
	if !exists {
		return nil
	}

	if err := process.cmd.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("failed to stop plugin: %w", err)
	}

	delete(m.plugins, pluginName)
	return nil
}
