package utils

import "math"

// CalculateExpectedScore determines the probability of a player winning against an opponent
// based on their respective Elo ratings.
func CalculateExpectedScore(playerRating float64, opponentRating float64) float64 {
	return 1 / (1 + math.Pow(10, (opponentRating-playerRating)/400))
}

// CalculateNewRating computes a player's new Elo rating after a match.
// The 'win' parameter should be 1.0 for a win, 0.5 for a draw, and 0.0 for a loss.
// The 'k' parameter is the K-factor, which determines the rating's sensitivity.
func CalculateNewRating(playerRating float64, opponentRating float64, score int, k int) float64 {
	expectedScore := CalculateExpectedScore(playerRating, opponentRating)
	return playerRating + float64(k)*(float64(score)-expectedScore)
}
