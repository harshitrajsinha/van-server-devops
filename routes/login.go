package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/harshitrajsinha/van-server-devops/models"
	"github.com/joho/godotenv"
)

func GenerateToken(username string) (string, error) {

	expiration := time.Now().Add(30 * time.Minute) // Expiration set as 30 minute

	claims := &jwt.StandardClaims{
		ExpiresAt: expiration.Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Load JWT key
	_ = godotenv.Load()
	jwtKeyString := os.Getenv("JWT_KEY")

	signedToken, err := token.SignedString([]byte(jwtKeyString))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code: http.StatusBadRequest, Message: "Invalid Request body for authorization",
		})
		log.Println("Invalid Request body for authorization")
		return
	}

	// Load username and password
	_ = godotenv.Load()
	validUsername := os.Getenv("AUTH_USER")
	validPassword := os.Getenv("AUTH_PASS")

	valid := (credentials.Username == validUsername && credentials.Password == validPassword)

	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code: http.StatusBadRequest, Message: "Incorrect username or password for authorization",
		})
		log.Println("Incorrect username or password for authorization")
		return
	}

	tokenString, err := GenerateToken(credentials.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code: http.StatusInternalServerError, Message: "Failed to generate token for authorization",
		})
		log.Println("Failed to generate token for authorization")
		return
	}
	response := make([]map[string]string, 0)
	response = append(response, map[string]string{"token": tokenString})
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Code    int                 `json:"code"`
		Message string              `json:"message"`
		Data    []map[string]string `json:"data"`
	}{
		Code: http.StatusCreated, Message: "Authorization token generated successfully. Valid for next 30mins", Data: response,
	})
	log.Println("Authorization token generated successfully")
}
