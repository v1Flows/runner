package plugin

import (
	"fmt"
	"sync"
)

type Registry struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
	}
}

func (r *Registry) Register(p Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	details := p.Details()
	if _, exists := r.plugins[details.Name]; exists {
		return fmt.Errorf("plugin %s already registered", details.Name)
	}

	if err := p.Validate(); err != nil {
		return fmt.Errorf("plugin %s validation failed: %w", details.Name, err)
	}

	r.plugins[details.Name] = p
	return nil
}

func (r *Registry) Get(name string) (Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.plugins[name]
	return p, exists
}

func (r *Registry) List() []Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}
	return plugins
}
