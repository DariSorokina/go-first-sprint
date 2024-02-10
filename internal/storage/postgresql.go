package storage

import (
	"database/sql"
)

type postgresqlDB struct {
	db *sql.DB
}

func newPostgresqlDB(cofigString string) *postgresqlDB {
	db, err := sql.Open("pgx", cofigString)
	if err != nil {
		panic(err) // TODO
	}

	return &postgresqlDB{db: db}
}

func (postgresqlDB *postgresqlDB) close() {
	postgresqlDB.db.Close()
}
