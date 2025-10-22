package models

// map from player ID to their Elo Rating
type EloRatings map[int]float64

type Store interface {
	GetAchievements() ([]Achievement, error)
	GetHeadToHead(p1 int, p2 int) (HeadToHead, error)
	GetIndexPageData() (IndexPageData, error)
	GetPlayerEloRatings(ids [2]int) (EloRatings, error)
	GetPlayerGames(id int, limit int) ([]Game, error)
	GetPlayerProfile(id int) (PlayerProfile, error)
	InsertGameResult(r GameResult) (int64, error)
	InsertPlayerAchievements(id int, achievementIDs []AchievementID) error
	InsertPlayer(name string) (int64, error)
	UpdateEloRatings(players EloRatings) error
}
