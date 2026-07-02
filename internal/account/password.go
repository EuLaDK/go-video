package account

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const passwordHashIterations = 60000
const passwordHashScheme = "nextvideo_hmac_sha256"

// hashPassword 生成带盐密码哈希；password 为用户提交的明文密码。
func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := derivePasswordHash([]byte(password), salt, passwordHashIterations)

	return strings.Join([]string{
		passwordHashScheme,
		strconv.Itoa(passwordHashIterations),
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	}, "$"), nil
}

// verifyPassword 校验明文密码；passwordHash 为 hashPassword 生成的持久化值。
func verifyPassword(password string, passwordHash string) bool {
	parts := strings.Split(passwordHash, "$")
	if len(parts) != 4 || parts[0] != passwordHashScheme {
		return false
	}

	iterations, err := strconv.Atoi(parts[1])
	if err != nil || iterations <= 0 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}
	wantHash, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false
	}

	gotHash := derivePasswordHash([]byte(password), salt, iterations)
	return subtle.ConstantTimeCompare(gotHash, wantHash) == 1
}

// derivePasswordHash 执行迭代 HMAC-SHA256；password/salt 为输入，iterations 为迭代次数。
func derivePasswordHash(password []byte, salt []byte, iterations int) []byte {
	current := append([]byte(nil), salt...)
	for index := 0; index < iterations; index++ {
		mac := hmac.New(sha256.New, password)
		mac.Write(current)
		current = mac.Sum(nil)
	}

	return current
}

// validatePassword 检查密码强度；password 为用户提交的明文密码。
func validatePassword(password string) error {
	if len(strings.TrimSpace(password)) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters", ErrInvalidAuthInput)
	}

	return nil
}

// isInvalidAuthInput 判断错误是否属于认证输入错误；err 为服务层错误。
func isInvalidAuthInput(err error) bool {
	return errors.Is(err, ErrInvalidAuthInput)
}
