package jwtutil

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"time"
)

var (
	accessTokenID        uint64
	refreshTokenID       uint64
	accessTokenIDPrefix  string
	refreshTokenIDPrefix string

	accessTokenPublicKey   = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	accessTokenPrivateKey  = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	accessTokenLifetime, _ = time.ParseDuration(os.Getenv("ACCESS_TOKEN_LIFETIME"))
	accessTokenAudiences   = []string{"web-scraper.dev"}

	// refreshTokenPublicKey   = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	refreshTokenPrivateKey  = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
	refreshTokenLifetime, _ = time.ParseDuration(os.Getenv("REFRESH_TOKEN_LIFETIME"))
	refreshTokenAudiences   = []string{"web-scraper.dev"}

	errInvalidToken = errors.New("invalid token")
)

func init() {
	var aBuf, rBuf [12]byte
	var aB64, rB64 string
	for len(aB64) < 10 {
		rand.Read(aBuf[:])
		aB64 = base64.StdEncoding.EncodeToString(aBuf[:])
		aB64 = strings.NewReplacer("+", "", "/", "").Replace(aB64)
	}
	accessTokenIDPrefix = aB64[0:10]

	for len(rB64) < 10 {
		rand.Read(rBuf[:])
		rB64 = base64.StdEncoding.EncodeToString(rBuf[:])
		rB64 = strings.NewReplacer("+", "", "/", "").Replace(rB64)
	}
	refreshTokenIDPrefix = rB64[0:10]
}
