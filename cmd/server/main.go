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
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/ws"
)

func main() {
	fmt.Println("Hello Backend Server!")
	config.Load()
	db.InitDB()

	userRepo := repositories.NewUserRepository(db.DB)
	authService := services.NewAuthService(userRepo)
	authHandler := api.NewAuthHandler(authService)

	chatroomHandler := api.NewChatroomHandler(services.NewChatroomService(
		*repositories.NewChatroomRepository(db.DB),
	))

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, authHandler)
	api.RegisterChatroomRoutes(mux, chatroomHandler)

	hub := ws.NewHub()
	go hub.Run()
	msgService := services.NewMessageService(repositories.NewMessageRepository(db.DB))
	api.SetupRoutes(mux, hub, msgService)

	log.Println("Server running on :3000")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
