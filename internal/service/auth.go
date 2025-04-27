package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sementrof/prod1/internal/models"
	"github.com/sirupsen/logrus"
)

func LoggerFactory() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	return logger
}

type Payload struct {
	jwt.RegisteredClaims
}

func GenerateJWTToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	payload := jwt.MapClaims{
		"sub:":  user.Id,
		"name":  user.Name,
		"email": user.Email,
		"exp":   expirationTime.Unix(),
	}
	if err := godotenv.Overload("env/.env"); err != nil {
		log.Printf("error loading env variables: %s", err)
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func IsExistsJWT(tokenString string, logger *logrus.Logger) (*Payload, error) {
	if err := godotenv.Load("env/.env"); err != nil {
		logger.Errorf("Error loading .env file: %v", err)
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}
	secretKey := os.Getenv("JWT_SECRET")

	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		logger.Errorf("Error parsing token: %v", err)
		return nil, err
	}
	claims, ok := token.Claims.(*Payload)
	if !ok || !token.Valid {
		logger.Errorf("Invalid token")
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		logger.Errorf("Token expired")
		return nil, fmt.Errorf("token expired")
	}
	return claims, nil
}

func Middleware(logger *logrus.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			logger.Error("Missing Authorization header")
			return
		}

		// Ожидаем токен в формате "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			logger.Error("Invalid token format")
			return
		}

		tokenString := parts[1]
		claims, err := IsExistsJWT(tokenString, logger)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			logger.Errorf("Invalid token: %v", err)
			return
		}

		// Можно использовать данные из токена
		logger.Infof("User ID from token: %s", claims.ID)

		// Переходим к следующему обработчику
		next(w, r)
	}
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You have accessed a protected route!"))
}
