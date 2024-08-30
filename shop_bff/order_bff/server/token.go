package server

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"order_bff/global"
	"time"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func GenJwtToken(payload string) (TokenPair, error) {

	expireTime := global.ServerConfig.Jwt.AccessExpire
	currentTime := time.Now().Unix()

	accessTokenClaims := jwt.MapClaims{
		"exp":     currentTime + expireTime,
		"iat":     currentTime,
		"payload": payload,
	}
	//使用jwt.NewWithClaims函数创建一个新的JWT令牌，指定了签名方法为HS256（HMAC SHA-256）和之前创建的声明。
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	accessSign, err := accessToken.SignedString([]byte(global.ServerConfig.Jwt.Key))
	if err != nil {
		return TokenPair{}, err
	}

	refreshTokenClaims := jwt.MapClaims{
		"exp": currentTime + expireTime,
		"iat": currentTime,
		"sub": payload,
	}
	//使用jwt.NewWithClaims函数创建一个新的JWT令牌，指定了签名方法为HS256（HMAC SHA-256）和之前创建的声明。
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	refreshSigning, err := refreshToken.SignedString([]byte(global.ServerConfig.Jwt.Key))
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessSign,
		RefreshToken: refreshSigning,
	}, nil
}
func CheckJwtToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.ServerConfig.Jwt.Key), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if payload, ok := claims["payload"].(string); ok {
			return payload, nil
		}
		return "", errors.New("无效的载荷")
	}

	return "", errors.New("无效的令牌")
}

// RefreshAccessToken 使用刷新令牌生成新的访问令牌
// @refreshTokenString: 刷新令牌
// @newAccessTokenExpireTime: 新的访问令牌过期时间，单位秒
func RefreshAccessToken(refreshTokenString string) (string, error) {
	newAccessTokenExpireTim := global.ServerConfig.Jwt.RefreshExpire
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.ServerConfig.Jwt.Key), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			currentTime := time.Now().Unix()

			newAccessTokenClaims := jwt.MapClaims{
				"exp":     currentTime + newAccessTokenExpireTim,
				"iat":     currentTime,
				"payload": sub,
			}
			newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessTokenClaims)
			return newAccessToken.SignedString([]byte(global.ServerConfig.Jwt.Key))
		}
		return "", errors.New("无效的载荷")
	}

	return "", errors.New("无效的令牌")
}
