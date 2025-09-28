package models

import "time"

type LeaderboardRow struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	EloRating float64 `json:"eloRating"`
}

// ---------------------------------------- achievements

type Achievement struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AchievementID int

// ---------------------------------------- players

type PlayerID struct {
	ID int `json:"id" binding:"required,min=1"`
}

type Name struct {
	Name string `json:"name" binding:"required,min=1,max=63"`
}

type Player struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PlayerProfile struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	EloRating    float64       `json:"eloRating"`
	CreatedAt    time.Time     `json:"createdAt"`
	GamesPlayed  int           `json:"gamesPlayed"`
	GamesWon     int           `json:"gamesWon"`
	RecentGames  []Game        `json:"recentGames"`
	Achievements []Achievement `json:"achievements"`
}

// ---------------------------------------- games

// Pointer values encode as the value pointed to. A nil pointer encodes as the null JSON object.
// Store a pointer to an int which will be encoded as an int if not nil and will be encoded as "null" if nil.

// Also see: https://www.sohamkamani.com/golang/omitempty/
type GameResult struct {
	WinnerID    int  `json:"winnerId" binding:"required"`
	LoserID     int  `json:"loserId" binding:"required"`
	WinnerScore *int `json:"winnerScore,omitempty" binding:"omitempty,min=0,max=255"`
	LoserScore  *int `json:"loserScore,omitempty" binding:"omitempty,min=0,max=255"`
}

type Game struct {
	ID          int       `json:"id"`
	Winner      Player    `json:"winner"`
	Loser       Player    `json:"loser"`
	WinnerScore *int      `json:"winnerScore"`
	LoserScore  *int      `json:"loserScore"`
	CreatedAt   time.Time `json:"createdAt"`
}
