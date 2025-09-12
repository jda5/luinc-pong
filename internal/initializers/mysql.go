package stores

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jda5/table-tennis/internal/stores"
)

func CreateMySQLStore() *stores.Connection {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err)
	}
	return &stores.Connection{DB: db}
}
