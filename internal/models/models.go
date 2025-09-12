package models

import "time"

type LeaderboardRow struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	EloRating float64 `json:"eloRating"`
}

// ---------------------------------------- players

type Name struct {
	Name string `json:"name" binding:"required,min=1,max=63"`
}

type Player struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	EloRating float64   `json:"eloRating"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
