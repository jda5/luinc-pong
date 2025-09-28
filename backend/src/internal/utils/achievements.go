package utils

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/jda5/luinc-pong/src/internal/models"
	"github.com/jda5/luinc-pong/src/internal/stores"
)

// -------------------------------------------------------------------------------- constants & types

const LIMIT int = 100_000

const (
	PLAY_1   models.AchievementID = iota + 1 // warming up
	PLAY_10                                  // minimum viable pong
	PLAY_50                                  // regular
	PLAY_100                                 // centurion
	PLAY_250                                 // legend
	PLAY_500                                 // unicorn

	WIN_11_0                     // chocolate
	WIN_11_1                     // bottle job
	WIN_12_10                    // clutch
	WIN_WITH_MORE_THAN_14_POINTS // marathon madness
	LOSE_12_10                   // heartbreaker

	WIN_5_CONSECUTIVE  // streaky
	WIN_10_CONSECUTIVE // unstoppable
	WIN_15_CONSECUTIVE // immortal

	LOSE_5_CONSECUTIVE                            // i get knocked down, but i get up again
	DAILY_WIN_3_CONSECUTIVE_AGAINST_SAME_OPPONENT // hat trick
	DAILY_WIN_5_CONSECUTIVE_AGAINST_SAME_OPPONENT // brutal

	LOSE_OPPONENT_15 // nemesis
	PLAY_OPPONENT_25 // rivalry
	PLAY_5_OPPONENTS // social butterfly

	PLAY_5_DAY  // daily standup
	PLAY_10_DAY // do you even work here?

	PLAY_OUTSIDE_WORK_HOURS // go home
	PLAY_3_DAY_STREAK       // dedicated
	PLAY_5_DAY_STREAK       // addicted

	WIN_UPSET_100_ELO // hostile takeover
	ELO_REACH_1100    // rising star
	ELO_REACH_1200    // big shot
	ELO_REACH_1300    // final boss
)

type DayCount struct {
	date  time.Time
	count int
}

type HeadToHead struct {
	playCount    int
	loseCount    int
	dayWinStreak DayCount
}

type AchievementSet map[models.AchievementID]struct{}

func (a AchievementSet) InsertID(id models.AchievementID) {
	a[id] = struct{}{}
}

// -------------------------------------------------------------------------------- public functions

// Calculate and update the achievements for both players based on their game history and recent game result.
func UpdatePlayerAchievements(
	s stores.Store,
	lastGame models.GameResult,
	oldRatings stores.EloRatings,
	newRatings stores.EloRatings,
) error {

	for _, id := range []int{lastGame.WinnerID, lastGame.LoserID} {
		games, err := s.GetPlayerGames(id, LIMIT)
		if err != nil {
			return fmt.Errorf("error updating player achievements %v", err)
		}
		if len(games) == 0 {
			// nothing to update
			continue
		}

		playerAchievements, err := calculatePlayersAchievements(
			id, games, lastGame, oldRatings, newRatings,
		)
		if err != nil {
			return fmt.Errorf("error updating player achievements %v", err)
		}
		err = s.InsertPlayerAchievements(id, playerAchievements)
		if err != nil {
			return fmt.Errorf("error updating player achievements %v", err)
		}
	}
	return nil
}

// -------------------------------------------------------------------------------- private functions

// Calculate the achievements a player has earned based on their game history and recent game result.
func calculatePlayersAchievements(
	id int,
	playerGames []models.Game,
	lastGame models.GameResult,
	oldRatings stores.EloRatings,
	newRating stores.EloRatings,
) ([]models.AchievementID, error) {

	a := make(AchievementSet)
	gamesPlayed := len(playerGames)

	// ---------------------------------------- game count milestones
	if gamesPlayed >= 1 {
		a.InsertID(PLAY_1)
	}
	if gamesPlayed >= 10 {
		a.InsertID(PLAY_10)
	}
	if gamesPlayed >= 50 {
		a.InsertID(PLAY_50)
	}
	if gamesPlayed >= 100 {
		a.InsertID(PLAY_100)
	}
	if gamesPlayed >= 250 {
		a.InsertID(PLAY_250)
	}
	if gamesPlayed >= 500 {
		a.InsertID(PLAY_500)
	}

	// ---------------------------------------- time-based achievements

	// A counter for consecutive days played
	playStreak := 0

	// A counter for consecutive wins
	winStreak := 0

	// A counter for consecutive losses
	loseStreak := 0

	// The last time the player played a game. Initialized to a date far in the past
	lastPlayed := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	var opponentID int
	var wonGame bool

	// A counter for number of games played in a single day
	dayStreak := 0

	// A map from opponent ID to head-to-head stats
	headToHeadMap := make(map[int]HeadToHead)

	for _, game := range playerGames {

		// ---------------------------------------- play streaks
		if datesEqual(game.CreatedAt, lastPlayed) {
			// same day as previous game
			dayStreak++
			switch dayStreak {
			case 5:
				a.InsertID(PLAY_5_DAY)
			case 10:
				a.InsertID(PLAY_10_DAY)
			}
		} else {
			// new day

			// day streak is reset
			dayStreak = 1

			if datesEqual(game.CreatedAt.AddDate(0, 0, -1), lastPlayed) {
				playStreak++
				switch playStreak {
				case 3:
					a.InsertID(PLAY_3_DAY_STREAK)
				case 5:
					a.InsertID(PLAY_5_DAY_STREAK)
				}
			} else {
				playStreak = 1
			}

		}
		lastPlayed = game.CreatedAt

		wonGame = game.Winner.ID == id
		if wonGame {
			opponentID = game.Loser.ID

			// ---------------------------------------- score-based
			if notNilPointer(game.WinnerScore) {
				if *game.WinnerScore == 11 {
					switch *game.LoserScore {
					case 0:
						a.InsertID(WIN_11_0)
					case 1:
						a.InsertID(WIN_11_1)
					}
				} else if *game.WinnerScore == 12 && *game.LoserScore == 10 {
					a.InsertID(WIN_12_10)
				} else if *game.WinnerScore >= 15 {
					a.InsertID(WIN_WITH_MORE_THAN_14_POINTS)
				}
			}

			// ---------------------------------------- winning streaks
			loseStreak = 0
			winStreak++
			switch winStreak {
			case 5:
				a.InsertID(WIN_5_CONSECUTIVE)
			case 10:
				a.InsertID(WIN_10_CONSECUTIVE)
			case 15:
				a.InsertID(WIN_15_CONSECUTIVE)
			}

			// ---------------------------------------- losing achievements
		} else {
			opponentID = game.Winner.ID
			winStreak = 0
			loseStreak++
			if loseStreak == 5 {
				a.InsertID(LOSE_5_CONSECUTIVE)
			}
			if notNilPointer(game.WinnerScore) && notNilPointer(game.LoserScore) {
				if *game.WinnerScore == 12 && *game.LoserScore == 10 {
					a.InsertID(LOSE_12_10)
				}
			}
		}

		// ---------------------------------------- head-to-head
		h, ok := headToHeadMap[opponentID]
		if !ok {
			headToHeadMap[opponentID] = HeadToHead{
				playCount: 1,
				loseCount: boolToInt(!wonGame),
				dayWinStreak: DayCount{
					date:  game.CreatedAt,
					count: boolToInt(game.Winner.ID == id),
				},
			}
		} else {
			h.playCount++
			if h.playCount == 25 {
				a.InsertID(PLAY_OPPONENT_25)
			}

			if wonGame {
				if datesEqual(h.dayWinStreak.date, game.CreatedAt) {
					h.dayWinStreak.count++
				} else {
					h.dayWinStreak.count = 1
					h.dayWinStreak.date = game.CreatedAt
				}

				switch h.dayWinStreak.count {
				case 3:
					a.InsertID(DAILY_WIN_3_CONSECUTIVE_AGAINST_SAME_OPPONENT)
				case 5:
					a.InsertID(DAILY_WIN_5_CONSECUTIVE_AGAINST_SAME_OPPONENT)
				}

			} else {
				h.dayWinStreak.count = 0
				h.dayWinStreak.date = game.CreatedAt
				h.loseCount++
				if h.loseCount == 15 {
					a.InsertID(LOSE_OPPONENT_15)
				}
			}

			// h is a copy of the struct value in the map; modifying it does not update the stored value
			headToHeadMap[opponentID] = h
		}

		// ---------------------------------------- misc
		playedAt := game.CreatedAt.Hour()
		if playedAt < 9 || playedAt > 17 {
			a.InsertID(PLAY_OUTSIDE_WORK_HOURS)
		}
	}

	if len(headToHeadMap) >= 5 {
		a.InsertID(PLAY_5_OPPONENTS)
	}

	achievements := make([]models.AchievementID, 0)
	if len(a) > 0 {
		achievements = slices.Collect(maps.Keys(a))
	}
	err := addPlayerEloAchievement(&achievements, id, lastGame, oldRatings, newRating)
	if err != nil {
		return achievements, err
	}
	return achievements, nil
}

func addPlayerEloAchievement(
	achievements *[]models.AchievementID,
	id int,
	lastGame models.GameResult,
	oldRatings stores.EloRatings,
	newRating stores.EloRatings,
) error {
	if id == lastGame.WinnerID {
		opponentElo, ok := oldRatings[lastGame.LoserID]
		if !ok {
			return fmt.Errorf("old elo rating for opponent not found")
		}
		playerElo, ok := oldRatings[id]
		if !ok {
			return fmt.Errorf("old elo rating for player not found")
		}
		if opponentElo-playerElo >= 100 {
			*achievements = append(*achievements, WIN_UPSET_100_ELO)
		}
	}

	newElo, ok := newRating[id]
	if !ok {
		return fmt.Errorf("error updating player achievements: old elo rating not found")
	}
	if newElo >= 1100 {
		*achievements = append(*achievements, ELO_REACH_1100)
	}
	if newElo >= 1200 {
		*achievements = append(*achievements, ELO_REACH_1200)
	}
	if newElo >= 1300 {
		*achievements = append(*achievements, ELO_REACH_1300)
	}
	return nil
}
