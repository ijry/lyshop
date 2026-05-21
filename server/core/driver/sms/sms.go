package sms

import (
	"context"
	"fmt"
	"sync"
)

// Driver is the interface all SMS plugins must implement.
type Driver interface {
	Name() string
	Send(ctx context.Context, phone, templateCode string, params map[string]string) error
}

var (
	mu     sync.RWMutex
	active Driver
)

func Register(d Driver) { mu.Lock(); active = d; mu.Unlock() }

func Get() (Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	if active == nil {
		return nil, fmt.Errorf("no SMS driver registered")
	}
	return active, nil
}
