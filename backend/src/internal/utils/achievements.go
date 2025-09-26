package utils

import (
	"github.com/jda5/luinc-pong/src/internal/models"
	"github.com/jda5/luinc-pong/src/internal/stores"
)

func UpdatePlayerAchievements(s stores.MySQLStore, r models.GameResult) {

}

// * First Blood: Play 1 game
// * Warm-up: Play 10 games
// * Play 50 games
// * Play 100 games
// * Play 250 games
// * Unicorn: Play 500 games
// * Play 5 games in a single day
// * Do You Even Work Here? Play 10 games in a single day

// * Chocolate: Win 11 - 0
// * Bottle Job: Win 11 - 1
// * Clutch Player: Win 12 - 10
// * Heartbreaker: Lose 10 - 12
// * Marathon Madness: Play a game that goes to 15+ points

// * Hostile Takeover: Beat some with 100 elo points more than you

// * Win 5 games in a row
// * Win 10 games in a row
// * Lose 5 games in a row
// * No Mercy: Win 3 games in a row against the same opponent in a day
// * Rival: Play the same opponent 25 times
// * Nemesis: Lose to someone 15 times
// * Social Player: Play 5 different people in the office.

// * Reach an ELO of 1100
// * Big Shot: Reach an ELO of 1200

// * Play a game on 3 consecutive days
// * Play a game on 5 consecutive days

// * Go Home: Play a game between the hours of 5 PM and 9 AM
