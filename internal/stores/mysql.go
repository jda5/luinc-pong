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
VALUES (?, ?, ?, ?)
`

var INSERT_PLAYER_QUERY string = `
INSERT INTO players (name)
VALUES (?)
`

var SELECT_LEADERBOARD_QUERY string = `
SELECT 
    id, name, elo_rating
FROM
    players
ORDER BY elo_rating DESC;
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

func (s *MySQLStore) InsertGameResult(r models.GameResult) (int64, error) {
	result, err := s.DB.Exec(INSERT_GAME_QUERY, r.WinnerID, r.LoserID, r.WinnerScore, r.LoserScore)
	if err != nil {
		return 0, fmt.Errorf("error inserting game: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error inserting game: %v", err)
	}
	return id, nil
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
