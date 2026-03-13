package models

import (
	"time"
)

// map from player ID to their Elo Rating
type EloRatings map[int]float64

type Store interface {
	DeleteGame(id int) error
	GetAchievements() ([]Achievement, error)
	GetGameResults() ([]BaseGame, error)
	GetGames(page int) ([]Game, error)
	GetHeadToHead(p1 int, p2 int) (HeadToHead, error)
	GetIndexPageData(showFull bool) (IndexPageData, error)
	GetPlayerBasicInfo() ([]PlayerBasicInfo, error)
	GetPlayerEloRatings(ids [2]int) (EloRatings, error)
	GetPlayerGames(id int, limit int) ([]Game, error)
	GetPlayerProfile(id int) (PlayerProfile, error)
	InsertGameResult(r GameResult) (int64, error)
	InsertPlayer(name string) (int64, error)
	InsertPlayerAchievements(id int, achievementIDs []AchievementID) error
	UpdateEloRatings(players EloRatings) error
	UpdateHighestEloRatings(players EloRatings) error
	UpdatePlayerUpdatedAt(m map[int]time.Time) error
}
