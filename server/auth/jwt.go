package auth

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/dgrijalva/jwt-go"
)

var (
	AccessTokenKey  = []byte("8f70563af940249929969adbc48fe7276354aa924934e28d5350e099ffec4028")
	ErrTokenInvalid = errors.New("jwt token invalid")
)

type claims struct {
	User   api.User `json:"user,omitempty"`
	UserId int32    `json:"user_id,omitempty"`
	jwt.StandardClaims
}

func ExpirationTime(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}

func CreateAccessToken(user *api.User, expirationTime time.Time, key []byte) (string, error) {
	claims := &claims{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.UnixNano(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func VerifyAccessToken(token string, key []byte) (*api.User, error) {
	claims := &claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, ErrTokenInvalid
	}

	return &claims.User, nil
}

func CreateRefreshToken(user *api.User, expirationTime time.Time, key []byte) (string, error) {
	claims := &claims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.UnixNano(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rToken, err := token.SignedString(key)
	return rToken, err
}

func VerifyRefreshToken(token string, key []byte) (int32, error) {
	claims := &claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		return key, nil
	})

	if err != nil {
		return 0, err
	}
	if !tkn.Valid {
		return 0, ErrTokenInvalid
	}

	return claims.UserId, nil
}

func CreateRandomKey() []byte {
	key := make([]byte, 32)
	rand.Read(key)

	return key
}
