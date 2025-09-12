package stores

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jda5/table-tennis/internal/stores/base"
)

func CreateMySQLStore() *base.Driver {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err)
	}
	return &base.Driver{db}
}
