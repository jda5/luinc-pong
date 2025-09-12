package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jda5/table-tennis/internal/models"
	"github.com/jda5/table-tennis/internal/stores"
)

type APIHandler struct {
	stores.Store
}

func (h *APIHandler) GetLeaderboard(c *gin.Context) {
	leaderboard, err := h.Store.GetLeaderboard()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, leaderboard)
}

func (h *APIHandler) InsertPlayer(c *gin.Context) {
	var name models.Name
	err := c.BindJSON(&name)
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	id, err := h.Store.InsertPlayer(name.Name)
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"message": fmt.Sprintf("a player with name `%s` already exists", name.Name)},
		)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"id": id})
}

func (h *APIHandler) InsertGame(c *gin.Context) {
	var result models.GameResult
	err := c.BindJSON(&result)
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	id, err := h.Store.InsertGameResult(result)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"id": id})

}
