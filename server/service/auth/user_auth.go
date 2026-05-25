package auth

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ijry/lyshop/core/cache"
	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/model"
	"gorm.io/gorm"
)

const smsCodeTTL = 5 * time.Minute
const smsCodeKeyPrefix = "sms:code:"

// SendSMSCode generates a 6-digit code and stores it in Redis for 5 min.
// Returns the code so the caller can pass it to the SMS driver.
func SendSMSCode(ctx context.Context, phone string) (string, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1000000)) //nolint:gosec
	key := smsCodeKeyPrefix + phone
	return code, cache.Client.Set(ctx, key, code, smsCodeTTL).Err()
}

// VerifySMSCode checks code against the stored value and deletes it on success.
func VerifySMSCode(ctx context.Context, phone, code string) error {
	key := smsCodeKeyPrefix + phone
	stored, err := cache.Client.Get(ctx, key).Result()
	if err != nil {
		return errors.New("auth.codeExpired")
	}
	if stored != code {
		return errors.New("auth.codeInvalid")
	}
	cache.Client.Del(ctx, key)
	return nil
}

// SMSLogin verifies the code, finds or creates the user, and returns a JWT.
func SMSLogin(ctx context.Context, phone, code string) (string, error) {
	if err := VerifySMSCode(ctx, phone, code); err != nil {
		return "", err
	}

	var user model.User
	err := db.DB.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = model.User{
			Phone:    phone,
			Nickname: "User" + phone[len(phone)-4:],
			Status:   1,
		}
		if err2 := db.DB.WithContext(ctx).Create(&user).Error; err2 != nil {
			return "", err2
		}
	} else if err != nil {
		return "", err
	}

	if user.Status == 0 {
		return "", errors.New("auth.accountDisabled")
	}

	return middleware.GenerateToken(user.ID, "user", nil)
}

// DeleteUserAccount verifies the SMS code, then soft-deletes the user.
func DeleteUserAccount(ctx context.Context, userID uint64, phone, code string) error {
	if err := VerifySMSCode(ctx, phone, code); err != nil {
		return err
	}

	// Verify phone matches the user
	var user model.User
	if err := db.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
		return errors.New("auth.userNotFound")
	}
	// Allow wx_ phone prefix (WeChat login users don't have real phone)
	if user.Phone != phone && user.Phone != "wx_"+phone {
		return errors.New("auth.phoneMismatch")
	}

	// Soft delete: set status=0 and anonymize
	return db.DB.WithContext(ctx).Model(&user).Updates(map[string]any{
		"status":   0,
		"phone":    "deleted_" + fmt.Sprintf("%d", userID),
		"nickname": "Deleted User",
		"avatar":   "",
	}).Error
}
