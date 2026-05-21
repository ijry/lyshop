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
		return errors.New("验证码已过期")
	}
	if stored != code {
		return errors.New("验证码错误")
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
			Nickname: "用户" + phone[len(phone)-4:],
			Status:   1,
		}
		if err2 := db.DB.WithContext(ctx).Create(&user).Error; err2 != nil {
			return "", err2
		}
	} else if err != nil {
		return "", err
	}

	if user.Status == 0 {
		return "", errors.New("账号已被禁用")
	}

	return middleware.GenerateToken(user.ID, "user", nil)
}
