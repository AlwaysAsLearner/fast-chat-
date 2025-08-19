package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlwaysAsLearner/fast-chat/backend/internal/services"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/utils"
)

type AuthHandler struct {
	Service *services.AuthService
}

type userData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req userData

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.Service.Register(req.Email, req.Password); err != nil {
		http.Error(w, fmt.Sprintf("Error During Registration: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userData
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user, err := h.Service.Authenticate(req.Email, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication Failed: %s", err), http.StatusBadRequest)
		return
	}

	// JWT stuff
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})

	w.Header().Add("Bearer", token)
}
