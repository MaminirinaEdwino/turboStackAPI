package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	turbojwt "github.com/MaminirinaEdwino/turbostackapi/src/TurboJwt"
	"github.com/MaminirinaEdwino/turbostackapi/src/activation"
	"github.com/MaminirinaEdwino/turbostackapi/src/config"
)

const JwtSecret = "turbostack_super_secret_key"

type UserService struct{}

func (s *UserService) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	_, err := config.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", creds.Username, creds.Password)
	if err != nil {
		http.Error(w, "User already exists or error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (s *UserService) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	var userID int
	err := config.DB.QueryRow("SELECT id FROM users WHERE username = $1 AND password = $2", creds.Username, creds.Password).Scan(&userID)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Génération du token JWT via turbojwt (durée 24h)
	token, err := turbojwt.Encode(JwtSecret, map[string]any{
		"user_id":  userID,
		"username": creds.Username,
	}, 24, 0)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *UserService) SaveToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token string `json:"token"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	activation.SaveTokenToFile(body.Token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Token Saved"})
}

func (s *UserService) CheckToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activation.CheckActivationToken())
}
