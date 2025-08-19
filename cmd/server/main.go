package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/AlwaysAsLearner/fast-chat/backend/configs"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/api"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/db"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/repositories"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/services"
)

func main() {
	fmt.Println("Hello Backend Server!")
	config.Load()
	db.InitDB()

	userRepo := repositories.NewUserRepository(db.DB)
	authService := services.NewAuthService(userRepo)
	authHandler := api.NewAuthHandler(authService)

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, authHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
