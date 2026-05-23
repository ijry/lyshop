package admin

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/model"
	"golang.org/x/crypto/bcrypt"
)

func ListAdmins(ctx context.Context) ([]model.Admin, error) {
	var list []model.Admin
	err := db.DB.WithContext(ctx).Find(&list).Error
	return list, err
}

func GetAdmin(ctx context.Context, id uint64) (*model.Admin, error) {
	var a model.Admin
	err := db.DB.WithContext(ctx).First(&a, id).Error
	return &a, err
}

func CreateAdmin(ctx context.Context, username, password string, roleID uint64) (*model.Admin, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	a := &model.Admin{
		Username: username,
		Password: string(hash),
		RoleID:   roleID,
		Status:   1,
	}
	return a, db.DB.WithContext(ctx).Create(a).Error
}

func UpdateAdmin(ctx context.Context, id uint64, updates map[string]any) error {
	// If password is being updated, hash it
	if pw, ok := updates["password"]; ok {
		hash, err := bcrypt.GenerateFromPassword([]byte(pw.(string)), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password"] = string(hash)
	}
	return db.DB.WithContext(ctx).Model(&model.Admin{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteAdmin(ctx context.Context, id uint64) error {
	return db.DB.WithContext(ctx).Delete(&model.Admin{}, id).Error
}

func ListRoles(ctx context.Context) ([]model.Role, error) {
	var list []model.Role
	err := db.DB.WithContext(ctx).Find(&list).Error
	return list, err
}

func GetRole(ctx context.Context, id uint64) (*model.Role, error) {
	var r model.Role
	err := db.DB.WithContext(ctx).First(&r, id).Error
	return &r, err
}

func CreateRole(ctx context.Context, name string, perms []string) (*model.Role, error) {
	permsJSON, _ := json.Marshal(perms)
	r := &model.Role{Name: name, Permissions: permsJSON}
	return r, db.DB.WithContext(ctx).Create(r).Error
}

func UpdateRole(ctx context.Context, id uint64, name string, perms []string) error {
	permsJSON, _ := json.Marshal(perms)
	return db.DB.WithContext(ctx).Model(&model.Role{}).Where("id = ?", id).
		Updates(map[string]any{"name": name, "permissions": permsJSON}).Error
}

func DeleteRole(ctx context.Context, id uint64) error {
	// Check no admins are using this role
	var count int64
	db.DB.Model(&model.Admin{}).Where("role_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该角色下还有管理员，无法删除")
	}
	return db.DB.WithContext(ctx).Delete(&model.Role{}, id).Error
}

// EnsureSuperAdmin creates the default super admin if no admins exist.
func EnsureSuperAdmin() error {
	var count int64
	db.DB.Model(&model.Admin{}).Count(&count)
	if count > 0 {
		return nil
	}

	// Create super admin role with wildcard permission
	superPerms, _ := json.Marshal([]string{"*"})
	role := &model.Role{Name: "超级管理员", Permissions: superPerms}
	if err := db.DB.FirstOrCreate(role, model.Role{Name: "超级管理员"}).Error; err != nil {
		return err
	}

	// Create admin user
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := &model.Admin{
		Username: "admin",
		Password: string(hash),
		RoleID:   role.ID,
		Status:   1,
	}
	return db.DB.Where(model.Admin{Username: "admin"}).FirstOrCreate(admin).Error
}
