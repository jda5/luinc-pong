package stores

import "database/sql"

type Connection struct {
	DB *sql.DB
}
