package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jda5/table-tennis/internal/models"
)

func GetLeaderboard(c *gin.Context) {
	leaderboard := []models.LeaderboardRow{
		{ID: 1, Name: "Gumbo", EloRating: 1400.9},
		{ID: 2, Name: "Jellife", EloRating: 1200.9},
		{ID: 3, Name: "Owen", EloRating: 1000.3},
		{ID: 4, Name: "Old Man", EloRating: 322.3},
	}
	c.IndentedJSON(http.StatusOK, leaderboard)
}
