package api

import (
	"fmt"
	"net/http"

	"github.com/AlwaysAsLearner/fast-chat/backend/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, authHandler *AuthHandler) {
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	mux.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey)
		w.Write([]byte("Your user ID: " + fmt.Sprint(userID)))
	})))
}
