package api

import (
	"fmt"
	"net/http"

	"github.com/AlwaysAsLearner/fast-chat/backend/internal/api/middleware"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/services"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/ws"
)

func RegisterRoutes(mux *http.ServeMux, authHandler *AuthHandler) {
	fmt.Println("Registering routes...")
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	mux.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey)
		w.Write([]byte("Your user ID: " + fmt.Sprint(userID)))
	})))
}

func SetupRoutes(mux *http.ServeMux, hub *ws.Hub, msgService *services.MessageService) {
	mux.Handle("/ws", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.WSServe(hub, msgService, w, r)
	})))
}
