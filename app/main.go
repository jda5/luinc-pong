package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jda5/table-tennis/internal/handlers"
	"github.com/jda5/table-tennis/internal/stores"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	router := gin.Default()
	h := handlers.APIHandler{Store: stores.CreateMySQLDAO()}

	router.GET("/leaderboards", h.GetLeaderboard)
	router.GET("/players/:id", h.GetPlayerProfile)
	router.POST("/players", h.InsertPlayer)
	router.POST("/games", h.InsertGame)

	router.Run("localhost:8080")
}
