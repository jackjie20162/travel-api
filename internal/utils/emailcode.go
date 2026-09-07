package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const emailCodePrefix = "email_code:"

// GenerateEmailCode 生成 6 位随机数字验证码
func GenerateEmailCode() (string, error) {
	var sb string
	for i := 0; i < 6; i++ {
		d, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		sb += fmt.Sprintf("%d", d)
	}
	return sb, nil
}

// StoreEmailCode 将邮件验证码存入 Redis
func StoreEmailCode(rds *redis.Redis, email, code string, expire int) error {
	key := emailCodePrefix + email
	return rds.Setex(key, code, expire)
}

// VerifyEmailCode 校验邮件验证码，校验成功后删除
func VerifyEmailCode(rds *redis.Redis, email, code string) bool {
	if email == "" || code == "" {
		return false
	}
	key := emailCodePrefix + email
	val, err := rds.Get(key)
	if err != nil || val == "" {
		return false
	}
	// 校验成功后删除（一次性）
	_, _ = rds.Del(key)
	return val == code
}
