package jwtutil

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"web-scraper.dev/internal/utils/ctxutil"
)

const issuer = "https://web-scraper.dev"

type AccessTokenClaims struct {
	TokenType string `json:"tokenType"`
	UserEmail string `json:"userEmail"`
	*jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	TokenType string `json:"tokenType"`
	UserEmail string `json:"userEmail"`
	*jwt.RegisteredClaims
}

func newAccessTokenClaims(userID, userEmail string) *AccessTokenClaims {
	id := fmt.Sprintf("%s-%06d", accessTokenIDPrefix, atomic.AddUint64(&accessTokenID, 1))

	return &AccessTokenClaims{
		TokenType: "access",
		UserEmail: userEmail,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			Audience:  accessTokenAudiences,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        id,
		},
	}
}

func newRefreshTokenClaims(userID, userEmail string) *RefreshTokenClaims {
	id := fmt.Sprintf("%s-%06d", refreshTokenIDPrefix, atomic.AddUint64(&refreshTokenID, 1))

	return &RefreshTokenClaims{
		TokenType: "refresh",
		UserEmail: userEmail,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			Audience:  refreshTokenAudiences,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        id,
		},
	}
}

func ClaimsFromAccessToken(token string) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Header["typ"] != "JWT" || token.Header["alg"] != jwt.SigningMethodEdDSA.Alg() {
			return nil, errInvalidToken
		}

		return jwt.ParseEdPublicKeyFromPEM([]byte(accessTokenPublicKey))
	})
	if err != nil {
		return nil, err
	}

	switch {
	case claims.ExpiresAt == nil,
		claims.Issuer != "https://web-scraper.dev",
		claims.Subject == "",
		claims.IssuedAt == nil,
		claims.ID == "",
		claims.TokenType != "access",
		claims.UserEmail == "":

		return nil, errInvalidToken
	}

	return claims, nil
}

func (c *AccessTokenClaims) ToCtxUser() ctxutil.User {
	ctxUser := ctxutil.User{
		Email: c.UserEmail,
	}

	userId, err := uuid.Parse(c.RegisteredClaims.Subject)
	if err == nil {
		ctxUser.ID = &userId
	}

	return ctxUser
}
