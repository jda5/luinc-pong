package stores

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jda5/luinc-pong/src/internal/models"
)

// ---------------------------------------- queries

const INSERT_GAME_QUERY string = `
INSERT INTO games (winner_id, loser_id, winner_score, loser_score)
VALUES (?, ?, ?, ?);
`

const INSERT_PLAYER_QUERY string = `
INSERT INTO players (name)
VALUES (?);
`

const SELECT_ACHIEVEMENTS_QUERY string = `
SELECT
	id, title, description
FROM
	achievement
`

const SELECT_PLAYER_ACHIEVEMENTS_QUERY string = `
SELECT 
   a.id, a.title, a.description
FROM
    achievement a
		JOIN
	player_achievement pa ON pa.achievement_id = a.id
WHERE
    pa.player_id = ?
ORDER BY pa.created_at DESC;
`

const SELECT_PLAYER_PROFILE_QUERY string = `
SELECT 
    (SELECT 
            COUNT(*)
        FROM
            games
        WHERE
            winner_id = players.id) AS total_wins,
    (SELECT 
            COUNT(*)
        FROM
            games
        WHERE
            loser_id = players.id) AS total_lost,
    name,
    elo_rating,
    created_at
FROM
    players
WHERE
    id = ?;
`

const SELECT_PLAYER_GAMES string = `
SELECT
	g.id AS game_id,
    w.id AS winner_id,
    w.name AS winner_name,
    l.id AS loser_id,
    l.name AS loser_name,
    g.winner_score,
    g.loser_score,
    g.created_at
FROM
    games g
        LEFT JOIN
    players w ON g.winner_id = w.id
        LEFT JOIN
    players l ON g.loser_id = l.id
WHERE
    g.winner_id = ? OR g.loser_id = ?
ORDER BY g.created_at DESC
LIMIT ?;
`

const SELECT_LEADERBOARD_QUERY string = `
SELECT 
    id, name, elo_rating
FROM
    players
ORDER BY elo_rating DESC;
`

const SELECT_PLAYER_ELO_RATINGS string = `
SELECT
	id, elo_rating
FROM
	players
WHERE
	id IN (?, ?);
`

const UPDATE_ELO_RATING_QUERY string = `
UPDATE players 
SET 
    elo_rating = ?
WHERE
    id = ?;
`

// ---------------------------------------- interface implementation

type MySQLStore struct {
	DB *sql.DB
	TZ *time.Location
}

func (s *MySQLStore) GetAchievements() ([]models.Achievement, error) {
	achievements := make([]models.Achievement, 0)

	rows, err := s.DB.Query(SELECT_ACHIEVEMENTS_QUERY)
	if err != nil {
		return nil, fmt.Errorf("error fetching achievements: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Achievement
		if err := rows.Scan(&a.ID, &a.Title, &a.Description); err != nil {
			return nil, fmt.Errorf("error fetching achievements: %v", err)
		}
		achievements = append(achievements, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching achievements: %v", err)
	}

	return achievements, nil
}

func (s *MySQLStore) GetLeaderboard() ([]models.LeaderboardRow, error) {
	leaderboard := make([]models.LeaderboardRow, 0)

	rows, err := s.DB.Query(SELECT_LEADERBOARD_QUERY)
	if err != nil {
		return nil, fmt.Errorf("error fetching leaderboard: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var row models.LeaderboardRow
		if err := rows.Scan(&row.ID, &row.Name, &row.EloRating); err != nil {
			return nil, fmt.Errorf("error fetching leaderboard: %v", err)
		}
		leaderboard = append(leaderboard, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching leaderboard: %v", err)
	}

	return leaderboard, nil
}

func (s *MySQLStore) GetPlayerEloRatings(ids [2]int) (EloRatings, error) {

	// need to use the make() function when creating a map
	ratings := make(EloRatings)

	rows, err := s.DB.Query(SELECT_PLAYER_ELO_RATINGS, ids[0], ids[1])
	if err != nil {
		return nil, fmt.Errorf("error geting player elo ratings: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var eloRating float64
		err := rows.Scan(&id, &eloRating)
		if err != nil {
			return nil, fmt.Errorf("error geting player elo ratings: %v", err)
		}
		ratings[id] = eloRating
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error geting player elo ratings: %v", err)
	}

	// check if we found the right number of rows
	if len(ratings) != 2 {
		return nil, fmt.Errorf("expected 2 players, but query found %d", len(ratings))
	}

	return ratings, nil
}

func (s *MySQLStore) GetPlayerGames(id int, limit int) ([]models.Game, error) {
	games := make([]models.Game, 0)
	rows, err := s.DB.Query(SELECT_PLAYER_GAMES, id, id, limit)
	if err != nil {
		return games, fmt.Errorf("error fetching games: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var game models.Game
		var winner models.Player
		var loser models.Player
		if err := rows.Scan(&game.ID, &winner.ID, &winner.Name, &loser.ID, &loser.Name, &game.WinnerScore, &game.LoserScore, &game.CreatedAt); err != nil {
			return games, fmt.Errorf("error fetching games: %v", err)
		}
		game.Winner = winner
		game.Loser = loser
		game.CreatedAt = game.CreatedAt.In(s.TZ)
	}
	if err := rows.Err(); err != nil {
		return games, fmt.Errorf("error fetching games: %v", err)
	}
	return games, nil
}

func (s *MySQLStore) GetPlayerProfile(id int) (models.PlayerProfile, error) {
	var profile models.PlayerProfile
	var totalWins int
	var totalLost int

	// ---------------------------------------- basic profile info
	row := s.DB.QueryRow(SELECT_PLAYER_PROFILE_QUERY, id)
	if err := row.Scan(&totalWins, &totalLost, &profile.Name, &profile.EloRating, &profile.CreatedAt); err != nil {
		return profile, fmt.Errorf("error fetching profile: %v", err)
	}
	profile.ID = id
	profile.GamesWon = totalWins
	profile.GamesPlayed = totalWins + totalLost
	profile.CreatedAt = profile.CreatedAt.In(s.TZ)

	// ---------------------------------------- recent games

	recentGames, err := s.GetPlayerGames(id, 20)
	if err != nil {
		return profile, fmt.Errorf("error fetching profile (recent games): %v", err)
	}
	profile.RecentGames = recentGames

	// ---------------------------------------- achievements

	achievementRows, err := s.DB.Query(SELECT_PLAYER_ACHIEVEMENTS_QUERY, id)
	if err != nil {
		return profile, fmt.Errorf("error fetching profile (achievements): %v", err)
	}
	defer achievementRows.Close()

	for achievementRows.Next() {
		var achievement models.Achievement
		if err := achievementRows.Scan(&achievement.ID, &achievement.Title, &achievement.Description); err != nil {
			return profile, fmt.Errorf("error fetching profile (achievements): %v", err)
		}
		profile.Achievements = append(profile.Achievements, achievement)
	}
	if err := achievementRows.Err(); err != nil {
		return profile, fmt.Errorf("error fetching profile (achievements): %v", err)
	}

	if len(profile.Achievements) == 0 {
		profile.Achievements = make([]models.Achievement, 0)
	}

	return profile, nil
}

func (s *MySQLStore) InsertGameResult(r models.GameResult) (int64, error) {
	result, err := s.DB.Exec(INSERT_GAME_QUERY, r.WinnerID, r.LoserID, r.WinnerScore, r.LoserScore)
	if err != nil {
		return 0, fmt.Errorf("error inserting game: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error inserting game: unknown player ID")
	}
	return id, nil
}

func (s *MySQLStore) InsertPlayer(name string) (int64, error) {
	result, err := s.DB.Exec(INSERT_PLAYER_QUERY, name)
	if err != nil {
		return 0, fmt.Errorf("error inserting player: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error inserting player: %v", err)
	}
	return id, nil
}

func (s *MySQLStore) InsertPlayerAchievements(id int, achievementIDs []models.AchievementID) error {
	if len(achievementIDs) == 0 {
		return nil // No achievements to insert; avoid invalid query
	}

	insertQuery := "INSERT IGNORE INTO player_achievement (player_id, achievement_id) VALUES "
	vals := []any{}

	for _, achievementID := range achievementIDs {
		insertQuery += "(?, ?),"
		vals = append(vals, id, int(achievementID))
	}
	// trim the last ,
	insertQuery = insertQuery[0 : len(insertQuery)-1]

	// prepare the statement
	stmt, err := s.DB.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer stmt.Close()

	// format all vals at once
	_, err = stmt.Exec(vals...)
	if err != nil {
		return fmt.Errorf("error inserting player achievements: %v", err)
	}
	return nil
}

func (s *MySQLStore) UpdateEloRatings(players EloRatings) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("error updating Players Elo rating: %v", err)
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	// Prepare the statement once for repeated use.
	stmt, err := tx.Prepare(UPDATE_ELO_RATING_QUERY)
	if err != nil {
		return fmt.Errorf("error updating Players Elo rating: %v", err)
	}

	for id, eloRating := range players {
		_, err := stmt.Exec(eloRating, id)
		if err != nil {
			return fmt.Errorf("error updating Player %v Elo rating: %v", id, err)
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error updating Players Elo rating: %v", err)
	}

	return nil
}

// ---------------------------------------- initialiser

func CreateMySQLDAO() *MySQLStore {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("MYSQL_USER")
	cfg.Passwd = os.Getenv("MYSQL_PASSWORD")

	cfg.Net = "tcp"
	cfg.Addr = "host.docker.internal:3306"
	cfg.DBName = "table_tennis"

	cfg.ParseTime = true

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(fmt.Sprintf("unable to connect to mysq dsn '%v': %v", cfg.FormatDSN(), err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("ping failed to mysq dsn '%v': %v", cfg.FormatDSN(), err))
	}

	tz, err := time.LoadLocation("Europe/London")
	if err != nil {
		panic(fmt.Sprintf("error fetching timezone: %v", err))
	}

	return &MySQLStore{DB: db, TZ: tz}
}
