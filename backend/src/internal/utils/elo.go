package utils

import (
	"math"
	"time"

	"github.com/jda5/luinc-pong/src/internal/models"
)

// CalculateExpectedScore determines the probability of a player winning against an opponent
// based on their respective Elo ratings.
func CalculateExpectedScore(playerRating float64, opponentRating float64) float64 {
	return 1 / (1 + math.Pow(10, (opponentRating-playerRating)/400))
}

// CalculateNewRating computes a player's new Elo rating after a match.
// The 'score' parameter should be 1 for a win, and 0 for a loss.
// The 'k' parameter is the K-factor, which determines the rating's sensitivity.
func CalculateNewRating(playerRating float64, opponentRating float64, score int, k int) float64 {
	expectedScore := CalculateExpectedScore(playerRating, opponentRating)
	return playerRating + float64(k)*(float64(score)-expectedScore)
}

func RecalculateEloRatings(s models.Store) error {

	// initialize all player ratings to 1000
	players, err := s.GetPlayerBasicInfo()
	if err != nil {
		return err
	}

	ratingMap := make(map[int]float64)
	lastPlayedMap := make(map[int]time.Time)

	for _, player := range players {
		ratingMap[player.ID] = 1000
		lastPlayedMap[player.ID] = player.CreatedAt
	}

	// fetch all games in chronological order
	games, err := s.GetGameResults()
	if err != nil {
		return err
	}

	for _, game := range games {
		// Get current ratings (default to 1000 if new player)
		winnerRating, ok := ratingMap[game.WinnerID]
		if !ok {
			winnerRating = 1000
		}
		loserRating, ok := ratingMap[game.LoserID]
		if !ok {
			loserRating = 1000
		}

		// Calculate and store new ratings
		ratingMap[game.WinnerID] = CalculateNewRating(winnerRating, loserRating, 1, 40)
		ratingMap[game.LoserID] = CalculateNewRating(loserRating, winnerRating, 0, 40)

		// Track last played time
		lastPlayedMap[game.WinnerID] = game.CreatedAt
		lastPlayedMap[game.LoserID] = game.CreatedAt
	}

	err = s.UpdateEloRatings(ratingMap)
	if err != nil {
		return err
	}

	err = s.UpdatePlayerUpdatedAt(lastPlayedMap)
	if err != nil {
		return err
	}

	return nil
}

// Pass interfaces by value. The interface itself is a small value, but it
// contains a pointer to the underlying concrete data (like *MySQLStore),
// so methods will correctly operate on the shared store instance.

// Returns the players' old ratings, their new ratings and an error.
func UpdatePlayersEloRating(s models.Store, winnerId int, loserId int) (models.EloRatings, models.EloRatings, error) {

	// fetch player elo rating
	ratingMap, err := s.GetPlayerEloRatings([2]int{winnerId, loserId})
	if err != nil {
		return ratingMap, ratingMap, err
	}

	winnerRating := ratingMap[winnerId]
	loserRating := ratingMap[loserId]

	// calculate new ratings
	ratingMap[winnerId] = CalculateNewRating(winnerRating, loserRating, 1, 40)
	ratingMap[loserId] = CalculateNewRating(loserRating, winnerRating, 0, 40)

	// update the ratings
	err = s.UpdateEloRatings(ratingMap)
	if err != nil {
		return ratingMap, ratingMap, err
	}

	oldRatings := models.EloRatings{
		winnerId: winnerRating,
		loserId:  loserRating,
	}
	return oldRatings, ratingMap, nil
}

func GetWinProbabilities(s models.Store, winnerId int, loserId int) (float64, float64, error) {
	ratingMap, err := s.GetPlayerEloRatings([2]int{winnerId, loserId})
	if err != nil {
		return 0, 0, err
	}

	winnerRating := ratingMap[winnerId]
	loserRating := ratingMap[loserId]

	winnerProb := CalculateExpectedScore(winnerRating, loserRating)
	loserProb := CalculateExpectedScore(loserRating, winnerRating)

	return winnerProb, loserProb, nil
}
