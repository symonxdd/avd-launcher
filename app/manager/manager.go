package manager

import (
	"avd-launcher/app/models"
	"context"
	"sync"
)

type AvdManager struct {
	ctx         context.Context
	runningAVDs map[string]*models.AVD
	mu          sync.RWMutex
}

func NewAvdManager() *AvdManager {
	return &AvdManager{
		runningAVDs: make(map[string]*models.AVD),
	}
}

func (m *AvdManager) setContext(ctx context.Context) {
	m.ctx = ctx
}

// SetContext provides a way for other packages to inject the context without exposing it to Wails bindings.
func SetContext(m *AvdManager, ctx context.Context) {
	m.setContext(ctx)
}

// Helper methods to manage running AVD state safely
func (m *AvdManager) setRunningAvd(name string, avd *models.AVD) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.runningAVDs[name] = avd
}

func (m *AvdManager) deleteRunningAvd(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.runningAVDs, name)
}

func (m *AvdManager) isAvdRunningInProcess(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.runningAVDs[name]
	return exists
}
