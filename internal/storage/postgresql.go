package storage

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresqlDB struct {
	db *sql.DB
}

func newPostgresqlDB(cofigBDString string) *postgresqlDB {
	db, err := sql.Open("pgx", cofigBDString)
	if err != nil {
		panic(err) // TODO
	}

	return &postgresqlDB{db: db}
}

func (postgresqlDB *postgresqlDB) close() {
	postgresqlDB.db.Close()
}
