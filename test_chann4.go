package main

import (
	"fmt"
	"time"
)

// Manager ...
type Manager struct {
	graceShut     chan struct{}
	triggerReload chan struct{}
}

// NewManager ...
func NewManager() *Manager {
	return &Manager{
		graceShut:     make(chan struct{}),
		triggerReload: make(chan struct{}, 1),
	}
}

func (m *Manager) reloader() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-m.graceShut:
			return
		case <-ticker.C:
			fmt.Println("ticker is firing...")
			// close(m.triggerReload)
			select {
			case _, ok := <-m.triggerReload:
				if ok {
					fmt.Println("trigger reload")
				} else {
					fmt.Println("trigger not reload")
				}
				// m.reload()
			case <-m.graceShut:
				return
			}
		}
	}
}

func main() {
	manager := NewManager()
	manager.reloader()
	// go manager.reloader()
}
