package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jda5/luinc-pong/src/internal/models"
	"github.com/jda5/luinc-pong/src/internal/utils"
)

type APIHandler struct {
	models.Store
}

// ---------------------------------------- internal helpers

func parsePlayerID(c *gin.Context, paramName string) (int, error) {
	idString := c.Query(paramName)
	if idString == "" {
		return 0, fmt.Errorf("missing required parameter: %s", paramName)
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, fmt.Errorf("'%s' is not a valid integer", idString)
	}

	if id <= 0 {
		return 0, fmt.Errorf("player ID must be positive")
	}

	return id, nil
}

// ---------------------------------------- public API

func (h *APIHandler) GetAchievements(c *gin.Context) {
	achievements, err := h.Store.GetAchievements()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, achievements)
}

func (h *APIHandler) GetIndexPage(c *gin.Context) {
	data, err := h.Store.GetIndexPageData()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, data)
}

func (h *APIHandler) GetPlayerProfile(c *gin.Context) {
	id, err := parsePlayerID(c, "id")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	profile, err := h.Store.GetPlayerProfile(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "player not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, profile)
}

func (h *APIHandler) GetHeadToHead(c *gin.Context) {
	p1, err := parsePlayerID(c, "p1")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	p2, err := parsePlayerID(c, "p2")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if p1 == p2 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "p1 and p2 must be different players"})
		return
	}

	headToHead, err := h.Store.GetHeadToHead(p1, p2)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, headToHead)
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

	if result.WinnerID == result.LoserID {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "winnerId and loserId cannot be the same"})
		return
	}

	id, err := h.Store.InsertGameResult(result)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	oldRatings, newRatings, err := utils.UpdatePlayersEloRating(h.Store, result.WinnerID, result.LoserID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	go func() {

		// Recover is a built-in function that regains control of a panicking goroutine.
		// Recover is only useful inside deferred functions.
		// 	// During normal execution, a call to recover will return nil and have no other effect.
		// 	// If the current goroutine is panicking, a call to recover will capture the value given
		// 	// to panic and resume normal execution.
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC recovered in UpdatePlayerAchievements: %v", r)
			}
		}()

		err := utils.UpdatePlayerAchievements(h.Store, result, oldRatings, newRatings)
		if err != nil {
			log.Printf("ERROR: background update of player achievements failed: %v", err)
		}
	}()

	c.IndentedJSON(http.StatusCreated, gin.H{"id": id})
}
