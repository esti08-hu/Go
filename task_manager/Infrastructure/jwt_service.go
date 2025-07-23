package infrastructure

import (
	"log"
	"os"
	"time"

	domain "task_manager/Domain"

	"github.com/joho/godotenv"
	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(user *domain.User) (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	var jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	
	jwtToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}