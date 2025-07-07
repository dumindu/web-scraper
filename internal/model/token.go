package model

import (
	"math/rand"
	"time"
)

const (
	tokenLength        = 6
	tokenLifetimeHours = 1
)

var alphaNum = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func NewToken() (string, *time.Time) {
	b := make([]rune, tokenLength)
	for i := range b {
		b[i] = alphaNum[rand.Intn(len(alphaNum))]
	}

	token := string(b)
	expiry := time.Now().Add(tokenLifetimeHours * time.Hour)
	return token, &expiry
}
