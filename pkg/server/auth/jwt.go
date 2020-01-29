package auth

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jheimbach/nfc-cash-system/api"
)

const (
	accessTokenKey   = "7QC/y4Dkke2izCGyArkfH074ETD9Hyf6PxIV/D7L2Nw="
	refreshTokenKey  = "tA2ZFqRCgYBEX4Y9/Q4Au9U0qrbW2oBcqJ8uRPavj9g="
	refreshTokenName = "refresh-tkn"
)

type TokenType string

const (
	AccessToken  TokenType = "access-tkn"
	RefreshToken TokenType = "refresh-tkn"
)

func (t TokenType) name() string {
	return string(t)
}

type claims struct {
	User   api.User `json:"user,omitempty"`
	UserId int32    `json:"user_id,omitempty"`
	jwt.StandardClaims
}

type TokenGenerator interface {
	ExpirationTime(duration time.Duration) time.Time
	CreateToken(user *api.User, expirationTime time.Time, tokenType TokenType) (string, error)
	VerifyToken(token string, tokenType TokenType) (user *api.User, expires time.Time, err error)
}

type JWTAuthenticator struct {
	keyStorage map[string][]byte
}

func NewJWTAuthenticator(accessTknKey, refreshTknKey string) (*JWTAuthenticator, error) {
	keyStorage := make(map[string][]byte)

	if accessTknKey == "" {
		return nil, fmt.Errorf("access token key must not be empty")
	}
	if refreshTknKey == "" {
		return nil, fmt.Errorf("refresh token key must not be empty")
	}
	keyStorage[string(AccessToken)] = []byte(accessTknKey)
	keyStorage[string(RefreshToken)] = []byte(refreshTknKey)

	return &JWTAuthenticator{keyStorage: keyStorage}, nil
}

func (JWTAuthenticator) ExpirationTime(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}

func (j JWTAuthenticator) CreateToken(user *api.User, expirationTime time.Time, tokenType TokenType) (string, error) {
	tokenName := tokenType.name()
	idHash := md5.New()
	io.WriteString(idHash, user.Name)
	io.WriteString(idHash, user.Email)
	io.WriteString(idHash, tokenName)

	claims := &claims{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			Id:        fmt.Sprintf("%x", idHash.Sum(nil)),
			ExpiresAt: expirationTime.Unix(),
			Subject:   fmt.Sprintf("user_%s_%d", user.Name, user.Id),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["type"] = tokenName

	return token.SignedString(j.keyStorage[tokenName])
}

func (j JWTAuthenticator) VerifyToken(token string, tokenType TokenType) (user *api.User, expires time.Time, err error) {
	claims := &claims{}
	_, e := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		headerType := token.Header["type"]
		tokenName := tokenType.name()
		if headerType != tokenName {
			return nil, fmt.Errorf("token is not from type %s", tokenType.name())
		}
		key := j.keyStorage[tokenName]
		return key, nil
	})

	if e != nil {
		return nil, time.Unix(0, 0), e
	}

	return &claims.User, time.Unix(claims.ExpiresAt, 0), nil
}
