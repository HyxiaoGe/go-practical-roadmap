package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourname/01-web-api-template/internal/config"
	"github.com/yourname/01-web-api-template/pkg/logger"
)

// Claims JWT声明结构体
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string) (string, error) {
	// 设置令牌过期时间
	expirationTime := time.Now().Add(time.Duration(config.GlobalConfig.JWT.AccessTokenExp) * time.Second)

	// 创建声明
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "access_token",
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证JWT令牌
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌有效性
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取Authorization字段
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Warn("Missing Authorization header")
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Warn("Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// 提取令牌
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 验证令牌
		claims, err := ValidateToken(tokenString)
		if err != nil {
			logger.Warn("Invalid token", logger.Error(err))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 将用户信息添加到请求上下文
		// 注意：这里简化处理，实际项目中可能需要更复杂的上下文管理

		// 调用下一个处理器
		next.ServeHTTP(w, r)
	})
}