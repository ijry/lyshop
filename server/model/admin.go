package model

// Admin is a back-office user.
type Admin struct {
	Base
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null"            json:"-"` // bcrypt hash
	RoleID   uint64 `gorm:"not null"                     json:"role_id"`
	Status   int8   `gorm:"not null;default:1"           json:"status"`
}
