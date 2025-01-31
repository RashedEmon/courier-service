package handlers

import (
	"courier-service/internal/database"
	"courier-service/internal/helpers"
	"courier-service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var userRequest models.SignUpRequest
	var user models.User

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		validationErrors := helpers.FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Please fix the given errors",
			"type":    "error",
			"code":    422,
			"errors":  validationErrors,
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Email = userRequest.Email
	user.Password = string(hashedPassword)

	if err := helpers.StoreUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

}

func Login(c *gin.Context) {
	var req models.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The user credentials were incorrect.",
			"type":    "error",
			"code":    400,
		},
		)
		return
	}

	// Generate JWT token
	accessToken, refreshToken, err := helpers.GenerateTokens(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Return token
	c.JSON(http.StatusOK, gin.H{
		"token_type":    "Bearer",
		"expires_in":    8 * 60 * 60,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	},
	)

}
