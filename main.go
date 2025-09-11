package main

import (
	"fmt"
	"internal/utils"
)

func main() {
	playerRating := 1000
	opponentRating := 1000
	score := 1
	k := 40
	newPlayerRating := utils.CalculateNewRating(
		playerRating, opponentRating, score, k,
	)
	fmt.Println(newPlayerRating)

	newOpponentRating := utils.CalculateNewRating(
		opponentRating, playerRating, 1.0-score, k,
	)
	fmt.Println(newOpponentRating)
}
