package auth

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminLogin validates username/password and returns a JWT with permissions.
func AdminLogin(ctx context.Context, username, password string) (string, error) {
	var admin model.Admin
	err := db.DB.WithContext(ctx).Where("username = ?", username).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", errors.New("auth.credentialsInvalid")
	}
	if err != nil {
		return "", err
	}
	if admin.Status == 0 {
		return "", errors.New("auth.accountDisabled")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return "", errors.New("auth.credentialsInvalid")
	}

	var role model.Role
	if err = db.DB.WithContext(ctx).First(&role, admin.RoleID).Error; err != nil {
		return "", err
	}
	var perms []string
	json.Unmarshal(role.Permissions, &perms) //nolint:errcheck

	return middleware.GenerateToken(admin.ID, "admin", perms)
}

// HashPassword returns the bcrypt hash of plain.
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}
