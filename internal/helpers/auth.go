package helpers

import (
	"courier-service/config"
	"courier-service/internal/database"
	"courier-service/internal/models"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func StoreUser(user models.User) error {
	if err := database.DB.Create(&user).Error; err != nil {

		return err
	}
	return nil
}

type Claims struct {
	Email  string `json:"email"`
	UserId uint   `json:"user_id"`
	jwt.StandardClaims
}

func GenerateTokens(email string, userId uint) (accessToken string, refreshToken string, err error) {
	// Define expiration times
	accessTokenExpiry := time.Now().Add(8 * time.Hour)       // 8 hours
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour) // 7 days

	// Create access token claims
	accessClaims := &Claims{
		Email:  email,
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpiry.Unix(),
		},
	}

	// Create refresh token claims
	refreshClaims := &Claims{
		Email:  email,
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiry.Unix(),
		},
	}

	// Generate tokens
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// Sign tokens with secret key
	accessToken, err = accessTokenObj.SignedString([]byte(config.ConfigInstance.JwtSecretKey))
	if err != nil {
		return "", "", err
	}

	refreshToken, err = refreshTokenObj.SignedString([]byte(config.ConfigInstance.JwtSecretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure token is signed with correct algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(config.ConfigInstance.JwtSecretKey), nil
	})

	if err != nil {
		return nil, err // Invalid token (signature mismatch, malformed, etc.)
	}

	// Extract claims and validate
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if the token has expired
		if claims.ExpiresAt < time.Now().Unix() {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
