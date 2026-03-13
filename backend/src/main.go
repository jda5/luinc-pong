package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jda5/luinc-pong/src/internal/handlers"
	"github.com/jda5/luinc-pong/src/internal/stores"

	"github.com/gin-contrib/cors"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	router := gin.Default()
	router.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)
	h := handlers.APIHandler{Store: stores.CreateMySQLDAO()}

	router.GET("/", h.GetIndexPage)
	router.GET("/achievements", h.GetAchievements)
	router.GET("/players/:id", h.GetPlayerProfile)
	router.GET("/head-to-head", h.GetHeadToHead)
	router.POST("/players", h.InsertPlayer)
	router.GET("/games", h.GetGames)
	router.DELETE("/games/:id", h.DeleteGame)
	router.POST("/games", h.InsertGame)
	router.GET("/recalculate", h.RecalculateElo)

	router.Run(":8080")
}
