package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jda5/table-tennis/internal/models"
	"github.com/jda5/table-tennis/internal/stores"
	"github.com/jda5/table-tennis/internal/utils"
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, err := h.Store.InsertGameResult(result)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Wrap the goroutine in an function literal to log any errors that have occured.
	go func() {

		// Recover is a built-in function that regains control of a panicking goroutine.
		// Recover is only useful inside deferred functions.
		// During normal execution, a call to recover will return nil and have no other effect.
		// If the current goroutine is panicking, a call to recover will capture the value given
		// to panic and resume normal execution.
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC recovered in UpdatePlayersEloRating: %v", r)
			}
		}()

		err := utils.UpdatePlayersEloRating(h.Store, result.WinnerID, result.LoserID)
		if err != nil {
			log.Printf("ERROR: background update of elo rating failed: %v", err)
		}
	}()

	c.IndentedJSON(http.StatusCreated, gin.H{"id": id})
}
