package utils

import (
	"math"

	"github.com/jda5/table-tennis/internal/stores"
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

// Pass interfaces by value. The interface itself is a small value, but it
// contains a pointer to the underlying concrete data (like *MySQLStore),
// so methods will correctly operate on the shared store instance.
func UpdatePlayersEloRating(s stores.Store, winnerId int, loserId int) error {

	// fetch player elo rating
	ratingMap, err := s.GetPlayerEloRatings([2]int{winnerId, loserId})
	if err != nil {
		return err
	}

	winnerRating := ratingMap[winnerId]
	loserRating := ratingMap[loserId]

	// calculate new ratings
	ratingMap[winnerId] = CalculateNewRating(winnerRating, loserRating, 1, 40)
	ratingMap[loserId] = CalculateNewRating(loserRating, winnerRating, 0, 40)

	// update the ratings
	return s.UpdateEloRatings(ratingMap)
}
