package jwtutil

import "github.com/golang-jwt/jwt/v5"

func NewAccessTokenAndRefreshToken(userID, userEmail string) (string, string, error) {
	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, newAccessTokenClaims(userID, userEmail))
	jwtAccessTokenPrivateKey, err := jwt.ParseEdPrivateKeyFromPEM([]byte(accessTokenPrivateKey))
	if err != nil {
		return "", "", err
	}
	accessToken, err := jwtAccessToken.SignedString(jwtAccessTokenPrivateKey)
	if err != nil {
		return "", "", err
	}

	jwtRefreshToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, newRefreshTokenClaims(userID, userEmail))
	jwtRefreshTokenPrivateKey, err := jwt.ParseEdPrivateKeyFromPEM([]byte(refreshTokenPrivateKey))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwtRefreshToken.SignedString(jwtRefreshTokenPrivateKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
