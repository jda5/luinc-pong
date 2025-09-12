package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jda5/table-tennis/internal/handlers"
)

func main() {
	router := gin.Default()
	router.GET("/leaderboard", handlers.GetLeaderboard)
	router.Run("localhost:8080")
}
