package jwtx

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

type JWT struct {
	secret string
	issuer string
	expire time.Duration
}

func (j *JWT) GetJWTSecret() []byte {
	return []byte(j.secret)
}

func (j *JWT) GenerateToke(userid int64) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(j.expire)
	claims := Claims{
		UserID: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    j.issuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(j.GetJWTSecret())
	return token, err
}

func (j *JWT) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.GetJWTSecret(), nil
	})
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func NewJWTEngine(Secret, Issuer string, Expire time.Duration) *JWT {
	return &JWT{
		secret: Secret,
		issuer: Issuer,
		expire: Expire,
	}
}
