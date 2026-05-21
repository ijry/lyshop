package oauth

import (
	"context"
	"fmt"
	"sync"
)

// UserInfo is the normalized user profile returned by any OAuth provider.
type UserInfo struct {
	OpenID   string
	UnionID  string
	Nickname string
	Avatar   string
}

// Driver is the interface all OAuth login plugins must implement.
type Driver interface {
	Name() string
	GetAuthURL(state string) string
	HandleCallback(ctx context.Context, code string) (*UserInfo, error)
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)
}

var (
	mu      sync.RWMutex
	drivers = map[string]Driver{}
)

func Register(d Driver) { mu.Lock(); drivers[d.Name()] = d; mu.Unlock() }

func Get(name string) (Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	d, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("oauth driver %q not registered", name)
	}
	return d, nil
}
