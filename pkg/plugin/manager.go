package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/AlertFlow/runner/config"
	"github.com/hashicorp/go-plugin"
)

type Manager struct {
	pluginDir string
	config    config.Config
	clients   map[string]*plugin.Client
	mutex     sync.RWMutex
}

// getDefaultPluginDir returns the default plugin directory in user's home
func getDefaultPluginDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	return filepath.Join(usr.HomeDir, ".alertflow", "plugins"), nil
}

func NewManager(config config.Config) (*Manager, error) {
	pluginDir := config.PluginDir

	// If no plugin directory specified, use default
	if pluginDir == "" {
		var err error
		pluginDir, err = getDefaultPluginDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get default plugin directory: %w", err)
		}
	}

	// Create plugin directory with parents
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory %s: %w", pluginDir, err)
	}

	// Verify directory is writable
	testFile := filepath.Join(pluginDir, ".write_test")
	if err := os.WriteFile(testFile, []byte{}, 0644); err != nil {
		return nil, fmt.Errorf("plugin directory %s is not writable: %w", pluginDir, err)
	}
	os.Remove(testFile)

	return &Manager{
		pluginDir: pluginDir,
		config:    config,
		clients:   make(map[string]*plugin.Client),
	}, nil
}

// downloadFromGitHub clones or updates a plugin from GitHub
func (m *Manager) downloadFromGitHub(ctx context.Context, pc config.PluginConfig) error {
	pluginPath := filepath.Join(m.pluginDir, pc.Name)

	// Create plugin-specific directory
	if err := os.MkdirAll(pluginPath, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory %s: %w", pluginPath, err)
	}

	// Check if repository already exists
	if _, err := os.Stat(filepath.Join(pluginPath, ".git")); os.IsNotExist(err) {
		// Clone repository
		cmd := exec.Command("git", "clone", pc.Repository, pluginPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	} else {
		// Update existing repository
		cmd := exec.Command("git", "fetch", "origin")
		cmd.Dir = pluginPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to fetch updates: %w", err)
		}
	}

	// Checkout specific version
	cmd := exec.Command("git", "checkout", pc.Version)
	cmd.Dir = pluginPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout version %s: %w", pc.Version, err)
	}

	// Build plugin
	cmd = exec.Command("go", "build", "-o", filepath.Join(pluginPath, pc.Name))
	cmd.Dir = pluginPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build plugin: %w", err)
	}

	return nil
}

// InstallPlugin downloads and builds a plugin
func (m *Manager) InstallPlugin(ctx context.Context, name string) error {
	var pluginConfig config.PluginConfig
	for _, pc := range m.config.Plugins {
		if pc.Name == name {
			pluginConfig = pc
			break
		}
	}

	if pluginConfig.Name == "" {
		return fmt.Errorf("plugin %s not found in configuration", name)
	}

	// Kill existing plugin instance if running
	m.mutex.Lock()
	if client, exists := m.clients[name]; exists {
		client.Kill()
		delete(m.clients, name)
	}
	m.mutex.Unlock()

	return m.downloadFromGitHub(ctx, pluginConfig)
}

// LoadPlugin loads a plugin and returns its client
func (m *Manager) LoadPlugin(name string) (*plugin.Client, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if client, exists := m.clients[name]; exists {
		return client, nil
	}

	client, err := m.startPlugin(name)
	if err != nil {
		return nil, err
	}

	m.clients[name] = client
	return client, nil
}

// startPlugin initializes and starts a new plugin client
func (m *Manager) startPlugin(name string) (*plugin.Client, error) {
	pluginPath := filepath.Join(m.pluginDir, name)

	config := &plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			name: &GRPCPlugin{
				// You can optionally provide a default implementation here
				Impl: nil,
			},
		},
		Cmd: exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
	}

	return plugin.NewClient(config), nil
}

// Add cleanup method
func (m *Manager) Cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for name, client := range m.clients {
		client.Kill()
		delete(m.clients, name)
	}
}
