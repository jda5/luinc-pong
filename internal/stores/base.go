package stores

import "github.com/jda5/table-tennis/internal/models"

// map from player ID to their Elo Rating
type EloRatings map[int]float64

type Store interface {
	GetLeaderboard() ([]models.LeaderboardRow, error)
	GetPlayerEloRatings(ids [2]int) (EloRatings, error)
	InsertGameResult(r models.GameResult) (int64, error)
	InsertPlayer(name string) (int64, error)
	UpdateEloRatings(players EloRatings) error
}
