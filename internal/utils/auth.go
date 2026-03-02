package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenDuration   = 60 * time.Minute
	RefreshTokenDuration  = 24 * time.Hour
	AccessTokenExpiresIn  = time.Hour
	ErrorAuthHeaderFormat = "invalid authorization header format"
)

var emailRegex *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}

	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	return true
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func VerifyPassword(password, hashedPassword string) (bool, error) {
	match, _, err := argon2id.CheckHash(password, hashedPassword)
	return match, err
}

func MakeJWT(userID int64, tokenSecret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "substrack",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiresIn)),
		Subject:   strconv.Itoa(int(userID)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (int64, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, fmt.Errorf("token is invalid")
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, err
	}

	return int64(userID), nil
}

func GetJWTToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf(ErrorAuthHeaderFormat)
	}

	const jwtPrefix = "JWT "
	if len(authHeader) < len(jwtPrefix) {
		return "", fmt.Errorf(ErrorAuthHeaderFormat)
	}

	if authHeader[:len(jwtPrefix)] != jwtPrefix {
		return "", fmt.Errorf(ErrorAuthHeaderFormat)
	}

	token := authHeader[len(jwtPrefix):]
	if token == "" {
		return "", fmt.Errorf(ErrorAuthHeaderFormat)
	}

	return token, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf(ErrorAuthHeaderFormat)
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) {
		return "", fmt.Errorf(ErrorAuthHeaderFormat)
	}

	if authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", fmt.Errorf(ErrorAuthHeaderFormat)
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		return "", fmt.Errorf(ErrorAuthHeaderFormat)
	}

	return token, nil
}

func MakeRefreshToken() (string, error) {
	randomData := make([]byte, 32)
	_, err := rand.Read(randomData)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomData), nil
}

func HashRefreshToken(refreshToken string) string {
	hash := sha512.Sum512([]byte(refreshToken))
	return hex.EncodeToString(hash[:])
}
