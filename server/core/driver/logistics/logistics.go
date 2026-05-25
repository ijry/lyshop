package logistics

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Driver interface {
	Name() string
	Query(ctx context.Context, req QueryReq) (*TrackResult, error)
}

type QueryReq struct {
	CompanyCode string `json:"company_code"`
	TrackingNo  string `json:"tracking_no"`
	Phone       string `json:"phone"`
}

type TrackNode struct {
	Time       time.Time       `json:"time"`
	Location   string          `json:"location"`
	StatusCode string          `json:"status_code"`
	StatusText string          `json:"status_text"`
	RawPayload json.RawMessage `json:"raw_payload"`
}

type TrackResult struct {
	Provider   string      `json:"provider"`
	StatusCode string      `json:"status_code"`
	StatusText string      `json:"status_text"`
	SignedAt   *time.Time  `json:"signed_at,omitempty"`
	Nodes      []TrackNode `json:"nodes"`
}

var (
	mu             sync.RWMutex
	drivers        = map[string]Driver{}
	defaultPrimary = "kuaidi100"
	defaultBackup  = "kdniao"
)

func Register(d Driver) {
	if d == nil {
		return
	}
	key := normalizeName(d.Name())
	if key == "" {
		return
	}
	mu.Lock()
	drivers[key] = d
	mu.Unlock()
}

func GetByName(name string) (Driver, error) {
	key := normalizeName(name)
	mu.RLock()
	defer mu.RUnlock()
	driver, ok := drivers[key]
	if !ok {
		return nil, fmt.Errorf("logistics driver %q not registered", key)
	}
	return driver, nil
}

func SetDefaultDrivers(primary, backup string) {
	mu.Lock()
	defer mu.Unlock()
	if normalized := normalizeName(primary); normalized != "" {
		defaultPrimary = normalized
	}
	if normalized := normalizeName(backup); normalized != "" {
		defaultBackup = normalized
	}
}

func DefaultDrivers() (string, string) {
	mu.RLock()
	defer mu.RUnlock()
	return defaultPrimary, defaultBackup
}

func ResolveByPinnedOrFallback(pinned string) (Driver, string, error) {
	pinnedName := normalizeName(pinned)
	if pinnedName != "" {
		driver, err := GetByName(pinnedName)
		return driver, pinnedName, err
	}
	primary, backup := DefaultDrivers()
	primaryDriver, primaryErr := GetByName(primary)
	if primaryErr == nil {
		return primaryDriver, primary, nil
	}
	backupDriver, backupErr := GetByName(backup)
	if backupErr == nil {
		return backupDriver, backup, nil
	}
	return nil, "", fmt.Errorf("primary=%v backup=%v", primaryErr, backupErr)
}

func normalizeName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "kuaidi100", "logistics_kuaidi100":
		return "kuaidi100"
	case "kdniao", "logistics_kdniao":
		return "kdniao"
	default:
		return strings.ToLower(strings.TrimSpace(name))
	}
}
