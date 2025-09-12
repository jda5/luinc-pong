package stores

import "github.com/jda5/table-tennis/internal/models"

type Store interface {
	GetLeaderboard() ([]models.LeaderboardRow, error)
	InsertPlayer(name string) (int64, error)
	InsertGameResult(r models.GameResult) (int64, error)
}
