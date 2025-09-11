package utils

import "testing"

func TestCalculateNewRating(t *testing.T) {
	playerRating := 1000
	opponentRating := 1000
	k := 40
	score := 1

	newRatings := [][]int{
		{1020, 980},
		{1038, 962},
		{1054, 946},
		{1068, 932},
	}

	for _, newRating := range newRatings {
		newPlayerRating := CalculateNewRating(playerRating, opponentRating, score, k)
		newOpponentRating := CalculateNewRating(opponentRating, playerRating, 1-score, k)
		if newPlayerRating != newRating[0] {
			t.Errorf("New player rating is incorrect. Expected %v, got %v", newRating[0], newPlayerRating)
		}
		if newOpponentRating != newRating[1] {
			t.Errorf("New opponent rating is incorrect. Expected %v, got %v", newRating[1], newOpponentRating)
		}
		playerRating = newPlayerRating
		opponentRating = newOpponentRating
	}
}

func TestCalculateExpectedScoreEqualRatings(t *testing.T) {
	playerRating := 1000
	opponentRating := 1000
	expectedScore := CalculateExpectedScore(playerRating, opponentRating)
	if expectedScore != 0.5 {
		t.Errorf("Expected score is incorrect. Expected 0.5, got %v", expectedScore)
	}
}

func TestCalculateExpectedScoreSumToOne(t *testing.T) {
	playerRating := 1025
	opponentRating := 2417
	if CalculateExpectedScore(playerRating, opponentRating)+CalculateExpectedScore(opponentRating, playerRating) != 1 {
		t.Error("Expected scores do not sum to 1")
	}
}
