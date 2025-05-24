package security

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTProvider struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type Claims struct {
	UserID *string
	Role   *string
}

func NewJWTProvider(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration) *JWTProvider {
	return &JWTProvider{
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (p *JWTProvider) GenerateAccessToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(p.accessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(p.secretKey))
}

func (p *JWTProvider) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(p.refreshTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(p.secretKey))
}

func (p *JWTProvider) VerifyAndParseClaims(tokenStr string) (Claims, error) {
	jwtClaims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.secretKey), nil
	})

	if err != nil {
		return Claims{}, err
	}

	claims := Claims{}

	userID, ok := jwtClaims["user_id"].(string)
	if ok {
		claims.UserID = &userID
	}

	role, ok := jwtClaims["role"].(string)
	if ok {
		claims.Role = &role
	}

	return claims, nil
}
