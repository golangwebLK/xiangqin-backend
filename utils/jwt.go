package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewJWT 初始化JWT
func NewJWT(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *JWT {
	if publicKey == nil {
		publicKey = &privateKey.PublicKey
	}
	return &JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// Sign 加密过程
func (j *JWT) Sign(claims jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", nil
	}
	return tokenString, err
}

// Verify 解密过程
func (j *JWT) Verify(tokenString string) (jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return j.publicKey, nil
		},
		jwt.WithValidMethods([]string{"RS256"}),
	)
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}
	return claims, nil
}

func NewJWTFromKeyBytes(keyBytes []byte) (*JWT, error) {
	// 解码给定的 PEM 数据，将其转换为一个 *pem.Block 结构
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	// 将传入的字节块解析为 PKCS1 格式的 RSA 私钥
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return NewJWT(key, nil), nil
}
