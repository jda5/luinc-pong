package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jda5/table-tennis/src/internal/handlers"
	"github.com/jda5/table-tennis/src/internal/stores"

	"github.com/gin-contrib/cors"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	router := gin.Default()
	router.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"http://localhost:5173", "https://luincpong.com"},
				AllowMethods:     []string{"GET", "POST", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)
	h := handlers.APIHandler{Store: stores.CreateMySQLDAO()}

	router.GET("/leaderboard", h.GetLeaderboard)
	router.GET("/players/:id", h.GetPlayerProfile)
	router.GET("/players", h.GetPlayerProfile)
	router.POST("/players", h.InsertPlayer)
	router.POST("/games", h.InsertGame)

	router.Run(":8080")
}
