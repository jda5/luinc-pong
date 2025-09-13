package stores

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jda5/table-tennis/internal/models"
)

// ---------------------------------------- queries

var INSERT_GAME_QUERY string = `
INSERT INTO games (winner_id, loser_id, winner_score, loser_score)
VALUES (?, ?, ?, ?);
`

var INSERT_PLAYER_QUERY string = `
INSERT INTO players (name)
VALUES (?);
`

var SELECT_LEADERBOARD_QUERY string = `
SELECT 
    id, name, elo_rating
FROM
    players
ORDER BY elo_rating DESC;
`

var SELECT_PLAYER_ELO_RATINGS string = `
SELECT
	id, elo_rating
FROM
	players
WHERE
	id IN (?, ?);
`

var UPDATE_ELO_RATING_QUERY string = `
UPDATE players 
SET 
    elo_rating = ?
WHERE
    id = ?;
`

// ---------------------------------------- interface implementation

type MySQLStore struct {
	DB *sql.DB
}

func (s *MySQLStore) GetLeaderboard() ([]models.LeaderboardRow, error) {
	var leaderboard []models.LeaderboardRow

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
	cfg.Addr = os.Getenv("MYSQL_HOST")
	cfg.DBName = os.Getenv("MYSQL_DATABASE")

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &MySQLStore{DB: db}
}
