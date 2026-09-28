package pkg

import (
	"errors"
	"fmt"
	"time"

	"gostart/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId    string `json:"userId"`
	TokenType string `json:"tokenType"` // token 类型，如 access 或 refresh
	jwt.RegisteredClaims
}

// ErrJWTSecretNotConfigured 表示 jwt.secretKey 未配置。
//
// 不提供任何兜底默认密钥：一旦回退到内置密钥，签发与校验就退化成公开值，
// 任何人都能伪造 token。宁可启动期失败，也不能带病运行。
var ErrJWTSecretNotConfigured = errors.New("jwt.secretKey 未配置：请在配置或环境变量中设置强密钥")

// jwtSecret 每次调用时读取配置，而不是在 init 期缓存。
// 这样测试可以自行设置 config.Configs，也避免包初始化顺序带来的空值。
func jwtSecret() ([]byte, error) {
	if config.Configs == nil || config.Configs.JWT == nil {
		return nil, ErrJWTSecretNotConfigured
	}
	secret := config.Configs.JWT.SecretKey
	if secret == "" {
		return nil, ErrJWTSecretNotConfigured
	}
	return []byte(secret), nil
}

func GenerateToken(uid string) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}

	now := time.Now()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{UserId: uid, TokenType: "access", RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenExpire())),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}})

	token, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}
	// 原实现在此处 fmt.Println 完整 token：日志落盘后等同于泄露可用的凭据，故移除。
	return token, nil
}

// tokenExpire 返回 token 有效期：优先用配置的 jwt.expire（秒），
// 未配置时沿用原有的 30 天默认值。
func tokenExpire() time.Duration {
	if config.Configs != nil && config.Configs.JWT != nil && config.Configs.JWT.Expire > 0 {
		return time.Duration(config.Configs.JWT.Expire) * time.Second
	}
	return 30 * 24 * time.Hour
}

func ParseAndValidateToken(tokenStr string, expectedTokenType string) (*CustomClaims, error) {
	secret, err := jwtSecret()
	if err != nil {
		return nil, err
	}

	jwtToken, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		// 锁定签名算法：防止攻击者把 alg 改成 none 或换成非对称算法来绕过校验。
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非法签名算法 %v，仅接受 HMAC", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := jwtToken.Claims.(*CustomClaims)
	if !ok || !jwtToken.Valid {
		return nil, errors.New("解析 claims 失败")
	}
	// 校验 Token 类型（防止拿 Refresh Token 去调用需要 Access Token 的业务接口）
	if claims.TokenType != expectedTokenType {
		return nil, fmt.Errorf("token 类型错误，期望 %s,实际为 %s", expectedTokenType, claims.TokenType)
	}
	return claims, nil
}
