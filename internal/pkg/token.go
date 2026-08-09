package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId    string `json:"userId"`
	TokenType string `json:"tokenType"` // token 类型，如 access 或 refresh
	jwt.RegisteredClaims
}

func GenerateToken(uid string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{UserId: uid, TokenType: "access", RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}})

	token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		return token, err
	}
	fmt.Println("AuthMiddleware token: ", token)
	return token, nil
}

func ParseAndValidateToken(tokenStr string, expectedTokenType string) (*CustomClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (any, error) { return []byte("secret"), nil })
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
