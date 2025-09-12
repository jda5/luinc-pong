package stores

import "database/sql"

type Driver struct {
	DB *sql.DB
}
