package models

type LeaderboardRow struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	EloRating float64 `json:"eloRating"`
}
