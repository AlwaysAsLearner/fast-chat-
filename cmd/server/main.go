package main

import (
	"fmt"

	config "github.com/AlwaysAsLearner/fast-chat/backend/configs"
	"github.com/AlwaysAsLearner/fast-chat/backend/internal/db"
)

func main() {
	fmt.Println("Hello Backend Server!")
	config.Load()
	fmt.Println("JWT config ", config.JWT)
	fmt.Println("Logger config ", config.Logger)
	fmt.Println("Websocket config ", config.WebSocket)
	fmt.Println("Database config ", config.DB)
	fmt.Println(config.DB.ToDsn())
	db.InitDB()
}
