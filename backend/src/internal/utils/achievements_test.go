package utils

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/jda5/luinc-pong/src/internal/models"
	"github.com/jda5/luinc-pong/src/internal/stores"
)

// -------------------------------------------------------------------------------- Test Helpers

func createTimeZone() *time.Location {
	tz, err := time.LoadLocation("Europe/London")
	if err != nil {
		panic(fmt.Sprintf("error creating timezone: %v", err))
	}
	return tz
}

// a helper function that creates a pointer to an interger
func intPointer(n int) *int {
	return &n
}

// gameToGameResult converts a models.Game to a models.GameResult for testing.
func gameToGameResult(g models.Game) models.GameResult {
	return models.GameResult{
		WinnerID:    g.Winner.ID,
		LoserID:     g.Loser.ID,
		WinnerScore: g.WinnerScore,
		LoserScore:  g.LoserScore,
	}
}

// createStandardEloRatings provides default ELO ratings for a game.
func createStandardEloRatings(g models.Game) (stores.EloRatings, stores.EloRatings) {
	return stores.EloRatings{
			g.Winner.ID: 1000,
			g.Loser.ID:  1000,
		}, stores.EloRatings{
			g.Winner.ID: 1016,
			g.Loser.ID:  984,
		}
}

// containsAchievement checks if a specific achievement is present in the results.
func containsAchievement(achievements []models.AchievementID, id models.AchievementID) bool {
	return slices.Contains(achievements, id)
}

// -------------------------------------------------------------------------------- Test Fixtures

var player models.Player = models.Player{
	ID:   1,
	Name: "Test Player",
}

var opponent models.Player = models.Player{
	ID:   2,
	Name: "Test Opponent",
}

var TZ *time.Location = createTimeZone()

// -------------------------------------------------------------------------------- Base Tests

func TestCalculatePlayersAchievementsReturnsPlayOne(t *testing.T) {
	playerGames := []models.Game{
		{
			ID:          1,
			Winner:      player,
			Loser:       opponent,
			WinnerScore: intPointer(11),
			LoserScore:  intPointer(5),
			CreatedAt:   time.Date(2023, 1, 2, 10, 0, 0, 0, TZ),
		},
	}
	lastGame := playerGames[len(playerGames)-1]
	oldRatings, newRating := createStandardEloRatings(lastGame)

	achievements, err := calculatePlayersAchievements(
		player.ID,
		playerGames,
		gameToGameResult(lastGame),
		oldRatings,
		newRating,
	)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !containsAchievement(achievements, PLAY_1) {
		t.Errorf("expected achievement PLAY_1, but it was not found")
	}
}

// -------------------------------------------------------------------------------- Added Tests

func TestWin15Consecutive(t *testing.T) {
	var playerGames []models.Game
	for i := range 15 {
		playerGames = append(playerGames, models.Game{
			ID:          i + 1,
			Winner:      player,
			Loser:       opponent,
			WinnerScore: intPointer(11),
			LoserScore:  intPointer(5),
			CreatedAt:   time.Date(2023, 1, 1, 10+i, 0, 0, 0, TZ),
		})
	}

	lastGame := playerGames[len(playerGames)-1]
	oldRatings, newRating := createStandardEloRatings(lastGame)

	achievements, err := calculatePlayersAchievements(player.ID, playerGames, gameToGameResult(lastGame), oldRatings, newRating)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsAchievement(achievements, WIN_15_CONSECUTIVE) {
		t.Errorf("expected achievement WIN_15_CONSECUTIVE, but it was not found")
	}
}

func TestPlay500(t *testing.T) {
	var playerGames []models.Game
	for i := range 500 {
		playerGames = append(playerGames, models.Game{
			ID:          i + 1,
			Winner:      player,
			Loser:       opponent,
			WinnerScore: intPointer(11),
			LoserScore:  intPointer(5),
			CreatedAt:   time.Date(2023, 1, 1, 0, i, 0, 0, TZ),
		})
	}

	lastGame := playerGames[len(playerGames)-1]
	oldRatings, newRating := createStandardEloRatings(lastGame)

	achievements, err := calculatePlayersAchievements(player.ID, playerGames, gameToGameResult(lastGame), oldRatings, newRating)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsAchievement(achievements, PLAY_500) {
		t.Errorf("expected achievement PLAY_500, but it was not found")
	}
}

func TestDailyWin5ConsecutiveAgainstSameOpponent(t *testing.T) {
	var playerGames []models.Game
	gameDay := time.Date(2023, 5, 5, 0, 0, 0, 0, TZ)
	for i := range 5 {
		playerGames = append(playerGames, models.Game{
			ID:          i + 1,
			Winner:      player,
			Loser:       opponent,
			WinnerScore: intPointer(11),
			LoserScore:  intPointer(5),
			CreatedAt:   gameDay.Add(time.Hour * time.Duration(i)),
		})
	}

	lastGame := playerGames[len(playerGames)-1]
	oldRatings, newRating := createStandardEloRatings(lastGame)

	achievements, err := calculatePlayersAchievements(player.ID, playerGames, gameToGameResult(lastGame), oldRatings, newRating)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsAchievement(achievements, DAILY_WIN_3_CONSECUTIVE_AGAINST_SAME_OPPONENT) {
		t.Errorf("expected achievement DAILY_WIN_3_CONSECUTIVE_AGAINST_SAME_OPPONENT, but it was not found")
	}

	if !containsAchievement(achievements, DAILY_WIN_5_CONSECUTIVE_AGAINST_SAME_OPPONENT) {
		t.Errorf("expected achievement DAILY_WIN_5_CONSECUTIVE_AGAINST_SAME_OPPONENT, but it was not found")
	}
}

func TestPlayOpponent25(t *testing.T) {
	var playerGames []models.Game
	d := time.Date(2023, 2, 1, 10, 0, 0, 0, TZ)

	for i := range 25 {
		// Alternate winner to ensure it's not win-streak dependent
		winner, loser := player, opponent
		if i%2 == 0 {
			winner, loser = opponent, player
		}
		playerGames = append(playerGames, models.Game{
			ID:          i + 1,
			Winner:      winner,
			Loser:       loser,
			WinnerScore: intPointer(11),
			LoserScore:  intPointer(5),
			CreatedAt:   d.Add(time.Hour * time.Duration(i)),
		})
	}

	lastGame := playerGames[len(playerGames)-1]
	oldRatings, newRating := createStandardEloRatings(lastGame)

	achievements, err := calculatePlayersAchievements(player.ID, playerGames, gameToGameResult(lastGame), oldRatings, newRating)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsAchievement(achievements, PLAY_OPPONENT_25) {
		t.Errorf("expected achievement PLAY_OPPONENT_25, but it was not found")
	}
}

func TestPlay5Opponents(t *testing.T) {
	var playerGames []models.Game
	for i := range 50 {
		opp := models.Player{ID: i + 2, Name: "Opponent"}
		playerGames = append(playerGames, models.Game{
			ID:          i + 1,
			Winner:      player,
			Loser:       opp,
			WinnerScore: intPointer(11),
			LoserScore:  intPointer(5),
			CreatedAt:   time.Date(2023, 3, 1+i, 10, 0, 0, 0, TZ),
		})
	}

	lastGame := playerGames[len(playerGames)-1]
	oldRatings, newRating := createStandardEloRatings(lastGame)

	achievements, err := calculatePlayersAchievements(player.ID, playerGames, gameToGameResult(lastGame), oldRatings, newRating)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsAchievement(achievements, PLAY_5_OPPONENTS) {
		t.Errorf("expected achievement PLAY_5_OPPONENTS, but it was not found")
	}
}

func TestPlay10Day(t *testing.T) {
	var playerGames []models.Game
	gameDay := time.Date(2023, 4, 1, 0, 0, 0, 0, TZ)
	for i := range 10 {
		playerGames = append(playerGames, models.Game{
			ID:          i + 1,
			Winner:      player,
			Loser:       opponent,
			WinnerScore: intPointer(11),
			LoserScore:  intPointer(5),
			CreatedAt:   gameDay.Add(time.Hour * time.Duration(i)),
		})
	}

	lastGame := playerGames[len(playerGames)-1]
	oldRatings, newRating := createStandardEloRatings(lastGame)

	achievements, err := calculatePlayersAchievements(player.ID, playerGames, gameToGameResult(lastGame), oldRatings, newRating)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsAchievement(achievements, PLAY_10_DAY) {
		t.Errorf("expected achievement PLAY_10_DAY, but it was not found")
	}
}

func TestWinUpset100Elo(t *testing.T) {
	playerGames := []models.Game{
		{
			ID:          1,
			Winner:      player, // The underdog winner
			Loser:       opponent,
			WinnerScore: intPointer(12),
			LoserScore:  intPointer(10),
			CreatedAt:   time.Date(2023, 6, 1, 10, 0, 0, 0, TZ),
		},
	}
	lastGame := playerGames[len(playerGames)-1]
	// Player ELO is 100 less than opponent's ELO
	oldRatings := stores.EloRatings{player.ID: 1000, opponent.ID: 1100}
	newRating := stores.EloRatings{player.ID: 1025, opponent.ID: 1075}

	achievements, err := calculatePlayersAchievements(player.ID, playerGames, gameToGameResult(lastGame), oldRatings, newRating)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsAchievement(achievements, WIN_UPSET_100_ELO) {
		t.Errorf("expected achievement WIN_UPSET_100_ELO, but it was not found")
	}
}

func TestEloReach1200(t *testing.T) {
	playerGames := []models.Game{
		{
			ID:          1,
			Winner:      player,
			Loser:       opponent,
			WinnerScore: intPointer(11),
			LoserScore:  intPointer(5),
			CreatedAt:   time.Date(2023, 7, 1, 10, 0, 0, 0, TZ),
		},
	}
	lastGame := playerGames[len(playerGames)-1]
	oldRatings := stores.EloRatings{player.ID: 1190, opponent.ID: 1050}
	// Player's new ELO crosses the 1200 threshold
	newRating := stores.EloRatings{player.ID: 1206, opponent.ID: 1034}

	achievements, err := calculatePlayersAchievements(player.ID, playerGames, gameToGameResult(lastGame), oldRatings, newRating)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsAchievement(achievements, ELO_REACH_1200) {
		t.Errorf("expected achievement ELO_REACH_1200, but it was not found")
	}
}

func TestPlay5DayStreak(t *testing.T) {
	var playerGames []models.Game
	startTime := time.Date(2023, 8, 1, 12, 0, 0, 0, TZ)
	for i := 0; i < 5; i++ {
		playerGames = append(playerGames, models.Game{
			ID:          i + 1,
			Winner:      player,
			Loser:       opponent,
			WinnerScore: intPointer(11),
			LoserScore:  intPointer(5),
			// Each game is on a consecutive day
			CreatedAt: startTime.AddDate(0, 0, i),
		})
	}

	lastGame := playerGames[len(playerGames)-1]
	oldRatings, newRating := createStandardEloRatings(lastGame)

	achievements, err := calculatePlayersAchievements(player.ID, playerGames, gameToGameResult(lastGame), oldRatings, newRating)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsAchievement(achievements, PLAY_5_DAY_STREAK) {
		t.Errorf("expected achievement PLAY_5_DAY_STREAK, but it was not found")
	}
}
