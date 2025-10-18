package stores

import "github.com/jda5/luinc-pong/src/internal/models"

// map from player ID to their Elo Rating
type EloRatings map[int]float64

type Store interface {
	GetAchievements() ([]models.Achievement, error)
	// GetLeaderboard() ([]models.LeaderboardRow, error)
	GetIndexPageData() (models.IndexPageData, error)
	GetPlayerEloRatings(ids [2]int) (EloRatings, error)
	GetPlayerGames(id int, limit int) ([]models.Game, error)
	GetPlayerProfile(id int) (models.PlayerProfile, error)
	InsertGameResult(r models.GameResult) (int64, error)
	InsertPlayerAchievements(id int, achievementIDs []models.AchievementID) error
	InsertPlayer(name string) (int64, error)
	UpdateEloRatings(players EloRatings) error
}
