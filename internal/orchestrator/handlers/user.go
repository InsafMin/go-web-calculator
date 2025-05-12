package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/InsafMin/go-web-calculator/internal/auth"
	"github.com/InsafMin/go-web-calculator/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE login = ?)", req.Login).Scan(&exists)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	_, err = db.DB.Exec("INSERT INTO users(login, password_hash) VALUES (?, ?)", req.Login, hashedPass)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "OK"}`)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user struct {
		ID           int
		PasswordHash string
	}
	row := db.DB.QueryRow("SELECT id, password_hash FROM users WHERE login = ?", req.Login)
	if err := row.Scan(&user.ID, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateToken(user.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
