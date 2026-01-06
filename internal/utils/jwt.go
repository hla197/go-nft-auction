package utils

import (
	"chain/internal/global"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// 密钥
	secretKey = global.JWT.Secret
	// 令牌有效期（访问令牌：1小时，刷新令牌：7天）
	tokenExpire = time.Hour * 24
)

type CustomClaims struct {
	UserID               uint64 `json:"user_id"`  // 用户ID
	Username             string `json:"username"` // 用户名
	jwt.RegisteredClaims        // 嵌入官方标准声明（包含exp/iss等）
}

// 生成JWT令牌（通用函数）
func GenerateToken(userID uint64, username string) (string, error) {
	// 构建自定义载荷
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建令牌（使用HS256算法）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名生成最终令牌字符串
	return token.SignedString([]byte(secretKey))
}

// 4. 验证并解析JWT令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{}, // 自定义载荷类型
		func(token *jwt.Token) (interface{}, error) {
			// 验证算法是否匹配
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		},
	)

	// 处理解析错误
	if err != nil {
		return nil, fmt.Errorf("parse token failed: %w", err)
	}

	// 验证令牌有效并提取载荷
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
