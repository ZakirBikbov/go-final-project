package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return nil
	}
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func generateToken(passwordHash string) (string, error) {
	secret := getJWTSecret()
	if secret == nil {
		return "", fmt.Errorf("password not configured")
	}

	claims := Claims{
		PasswordHash: passwordHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func validateToken(tokenString string) (string, error) {
	secret := getJWTSecret()
	if secret == nil {
		return "", fmt.Errorf("password not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.PasswordHash, nil
	}

	return "", fmt.Errorf("invalid token")
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if len(password) == 0 {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		tokenHash, err := validateToken(cookie.Value)
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		currentHash := hashPassword(password)
		if tokenHash != currentHash {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}

type signinRequest struct {
	Password string `json:"password"`
}

type signinResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, signinResponse{Error: "Неверный формат запроса"})
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		writeJSON(w, http.StatusBadRequest, signinResponse{Error: "Аутентификация не настроена"})
		return
	}

	if req.Password != password {
		writeJSON(w, http.StatusUnauthorized, signinResponse{Error: "Неверный пароль"})
		return
	}

	passwordHash := hashPassword(password)
	token, err := generateToken(passwordHash)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		writeJSON(w, http.StatusInternalServerError, signinResponse{Error: "Ошибка создания токена"})
		return
	}

	writeJSON(w, http.StatusOK, signinResponse{Token: token})
}
