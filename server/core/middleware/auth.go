package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ijry/lyshop/config"
	"github.com/ijry/lyshop/core/response"
)

// Claims is the JWT payload used for both user and admin tokens.
type Claims struct {
	UserID uint64   `json:"user_id"`
	Role   string   `json:"role"` // "user" | "admin"
	Perms  []string `json:"perms,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken signs a JWT for userID with the given role and permissions.
func GenerateToken(userID uint64, role string, perms []string) (string, error) {
	expiry := time.Duration(config.Global.JWT.ExpireHours) * time.Hour
	claims := Claims{
		UserID: userID,
		Role:   role,
		Perms:  perms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(config.Global.JWT.Secret))
}

// ParseToken validates tokenStr and returns the Claims.
func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims,
		func(_ *jwt.Token) (any, error) {
			return []byte(config.Global.JWT.Secret), nil
		},
	)
	return claims, err
}

// RequireAuth aborts with 401 if the request has no valid JWT.
func RequireAuth(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusOK, response.Err(401, "请先登录"))
		return
	}
	claims, err := ParseToken(strings.TrimPrefix(auth, "Bearer "))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, response.Err(401, "Token无效或已过期"))
		return
	}
	c.Set("user_id", claims.UserID)
	c.Set("role", claims.Role)
	c.Set("perms", claims.Perms)
	c.Next()
}

// RequireAdmin calls RequireAuth, then checks role == "admin".
func RequireAdmin(c *gin.Context) {
	RequireAuth(c)
	if c.IsAborted() {
		return
	}
	if role, _ := c.Get("role"); role != "admin" {
		c.AbortWithStatusJSON(http.StatusOK, response.Err(403, "无权限"))
		return
	}
	c.Next()
}

// RequirePermission returns a middleware that checks if the admin has the given permission.
// "*" in perms acts as a superadmin wildcard (all permissions granted).
func RequirePermission(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms, exists := c.Get("perms")
		if !exists {
			c.AbortWithStatusJSON(http.StatusOK, response.Err(403, "无权限"))
			return
		}
		permList, ok := perms.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusOK, response.Err(403, "无权限"))
			return
		}
		for _, p := range permList {
			if p == "*" || p == perm {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, response.Err(403, "无权限: "+perm))
	}
}

// HasPermission checks if perms slice contains perm (or wildcard "*").
func HasPermission(perms []string, perm string) bool {
	for _, p := range perms {
		if p == "*" || p == perm {
			return true
		}
	}
	return false
}
