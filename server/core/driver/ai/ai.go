package ai

import (
	"context"
	"fmt"
	"sync"
)

type GenerateParams struct {
	Prompt    string
	NegPrompt string
	Width     int
	Height    int
	Count     int    // 1-5
	Style     string // "ecommerce" | "realistic" | "illustration"
	RefImageURL string
}

type GenerateResult struct {
	URLs []string
}

// Driver is the interface all AI image generation drivers must implement.
type Driver interface {
	Name() string
	Generate(ctx context.Context, p *GenerateParams) (*GenerateResult, error)
}

var (
	mu      sync.RWMutex
	drivers = map[string]Driver{}
	def     string
)

func Register(d Driver, isDefault bool) {
	mu.Lock()
	defer mu.Unlock()
	drivers[d.Name()] = d
	if isDefault || def == "" {
		def = d.Name()
	}
}

func Get(name string) (Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	if name == "" {
		name = def
	}
	d, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("ai driver %q not registered", name)
	}
	return d, nil
}

func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(drivers))
	for n := range drivers {
		out = append(out, n)
	}
	return out
}
