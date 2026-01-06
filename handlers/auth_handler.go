package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/config"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	json.NewDecoder(r.Body).Decode(&creds)

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(creds.Password),
		bcrypt.DefaultCost,
	)

	err := repository.CreateUser(creds.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	json.NewDecoder(r.Body).Decode(&creds)

	user, storedHash, err := repository.GetUserByUsername(creds.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(config.JWTKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user":  user,
	})
}
