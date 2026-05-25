package delivery

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Driver is the interface every delivery method must implement.
type Driver interface {
	Name() string
	ValidateShipment(req ValidateReq) error
	NeedTrackingSync() bool
}

// ValidateReq carries the shipment parameters to validate.
type ValidateReq struct {
	Company    string // courier company code (express)
	TrackingNo string // tracking number (express)
	RiderName  string // rider name (local)
	RiderPhone string // rider phone (local)
}

var (
	mu      sync.RWMutex
	drivers = map[string]Driver{}
)

// Register adds a delivery driver to the global registry.
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

// Get returns the delivery driver by name.
func Get(name string) (Driver, error) {
	key := normalizeName(name)
	mu.RLock()
	defer mu.RUnlock()
	driver, ok := drivers[key]
	if !ok {
		return nil, fmt.Errorf("delivery driver %q not registered", key)
	}
	return driver, nil
}

// Names returns all registered delivery driver names.
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(drivers))
	for k := range drivers {
		names = append(names, k)
	}
	return names
}

// Validate validates shipment params using the specified delivery driver.
func Validate(driverName string, req ValidateReq) error {
	d, err := Get(driverName)
	if err != nil {
		return err
	}
	return d.ValidateShipment(req)
}

// NeedSync returns whether the delivery type requires logistics tracking sync.
func NeedSync(driverName string) bool {
	d, err := Get(driverName)
	if err != nil {
		return false
	}
	return d.NeedTrackingSync()
}

func normalizeName(name string) string {
	n := strings.ToLower(strings.TrimSpace(name))
	switch n {
	case "express", "delivery_express":
		return "express"
	case "local", "delivery_local":
		return "local"
	default:
		return n
	}
}

// Common validation errors.
var (
	ErrMissingCompany    = errors.New("delivery.err.expressRequired")
	ErrMissingTrackingNo = errors.New("delivery.err.trackingRequired")
	ErrMissingRiderName  = errors.New("delivery.err.riderNameRequired")
	ErrMissingRiderPhone = errors.New("delivery.err.riderPhoneRequired")
)
