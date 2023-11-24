package jwtAuth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenerateTokens(userId int64) (accessTokenString string, refreshTokenString string, err error) {
	const op = "lib.jwt.GenerateToken"

	jwtIssuer := os.Getenv("JWT_ISSUER")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	accessTokenClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    jwtIssuer,
		Subject:   strconv.FormatInt(userId, 10),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err = accessToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to sign access token: %w", op, err)
	}

	refreshTokenClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    jwtIssuer,
		Subject:   strconv.FormatInt(userId, 10),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err = refreshToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to sign refresh token: %w", op, err)
	}

	return accessTokenString, refreshTokenString, err
}

func ValidateToken(tokenString string) (*jwt.StandardClaims, error) {
	const op = "lib.jwt.ValidateToken"

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: unexpected signing method: %v", op, token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("%s: invalid token %w", op, err)
}

func ExtractToken(r *http.Request) string {
	const op = "lib.jwt.ExtractToken"

	bearToken := r.Header.Get("Authorization")

	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}