package auth

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/dgrijalva/jwt-go"
)

var (
	AccessTokenKey = []byte("8f70563af940249929969adbc48fe7276354aa924934e28d5350e099ffec4028")
)

type claims struct {
	User   api.User `json:"user,omitempty"`
	UserId int32    `json:"user_id,omitempty"`
	jwt.StandardClaims
}

type TokenGenerator interface {
	ExpirationTime(duration time.Duration) time.Time
	CreateToken(user *api.User, expirationTime time.Time, key []byte) (string, error)
	VerifyToken(token string, key []byte) (*api.User, error)
	CreateRandomKey() []byte
}

type jwtGenerator int

func NewJwtGenerator() TokenGenerator {
	return new(jwtGenerator)
}

func (*jwtGenerator) ExpirationTime(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}

func (*jwtGenerator) CreateToken(user *api.User, expirationTime time.Time, key []byte) (string, error) {
	idHash := md5.New()
	io.WriteString(idHash, user.Name)
	io.WriteString(idHash, string(key))

	claims := &claims{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			Id:        fmt.Sprintf("%x", idHash.Sum(nil)),
			ExpiresAt: expirationTime.Unix(),
			Subject:   fmt.Sprintf("user_%s_%d", user.Name, user.Id),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func (*jwtGenerator) VerifyToken(token string, key []byte) (*api.User, error) {
	claims := &claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	return &claims.User, nil
}

func (*jwtGenerator) CreateRandomKey() []byte {
	key := make([]byte, 32)
	rand.Read(key)

	return []byte(base64.StdEncoding.EncodeToString(key))
}
