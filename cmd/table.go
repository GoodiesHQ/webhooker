package main

import (
	"fmt"
	"sync"
)

type WebhookTable struct {
	webhooks map[string][]string // map to store the name and its list of targets
	mu       sync.RWMutex        // read/write mutex for modifying the webhooks
}

func (table *WebhookTable) Init() {
	// lock the mutex for writing
	table.mu.Lock()
	defer table.mu.Unlock()

	table.webhooks = make(map[string][]string)
}

func (table *WebhookTable) Set(name string, targets []string, override bool) error {
	// lock the mutex for writing
	table.mu.Lock()
	defer table.mu.Unlock()

	// check if the name exists
	_, found := table.webhooks[name]
	if found && !override {
		// if not told to override, it is a duplicate
		return fmt.Errorf("duplicate name '%s'", name)
	}

	table.webhooks[name] = targets
	return nil
}

func (table *WebhookTable) Get(name string) []string {
	// lock the mutex for reading
	table.mu.RLock()
	defer table.mu.RUnlock()

	// get the list of targets
	if targets, found := table.webhooks[name]; found {
		if targets == nil {
			return []string{} // return an empty list instead of nil
		}
		return targets
	}

	// name is not found
	return nil
}
