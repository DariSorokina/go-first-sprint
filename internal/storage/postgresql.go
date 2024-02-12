package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	createSchemaQuery = `CREATE SCHEMA IF NOT EXISTS content;`
	createTableQuery  = `CREATE TABLE IF NOT EXISTS content.urls (
		originalURL TEXT, 
		shortURL TEXT);`
	writeTestURLsQuery   = `INSERT INTO content.urls (originalURL, shortURL) VALUES ('https://practicum.yandex.ru/', 'd41d8cd98f');`
	readShortURLQuery    = `SELECT shortURL FROM content.urls WHERE originalURL = $1;`
	readOriginalURLQuery = `SELECT originalURL FROM content.urls WHERE shortURL = $1;`
	writeURLsQuery       = `INSERT INTO content.urls (originalURL, shortURL) VALUES ($1, $2);`
)

type PostgresqlDB struct {
	db *sql.DB
}

func NewPostgresqlDB(cofigBDString string) *PostgresqlDB {
	db, err := sql.Open("pgx", cofigBDString)
	if err != nil {
		panic(err) // TODO
	}

	return &PostgresqlDB{db: db}
}

func (postgresqlDB *PostgresqlDB) createTables() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error

	_, err = postgresqlDB.db.ExecContext(ctx, createSchemaQuery)
	if err != nil {
		return err
	}
	_, err = postgresqlDB.db.ExecContext(ctx, createTableQuery)
	if err != nil {
		return err
	}
	_, err = postgresqlDB.db.ExecContext(ctx, writeTestURLsQuery)
	if err != nil {
		return err
	}

	return nil
}

func (postgresqlDB *PostgresqlDB) SetValue(shortURL, longURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := postgresqlDB.db.ExecContext(ctx, writeURLsQuery, longURL, shortURL)
	if err != nil {
		log.Fatal(err)
	}

}

func (postgresqlDB *PostgresqlDB) GetShort(longURL string) (shortURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := postgresqlDB.db.QueryRowContext(ctx, readShortURLQuery, longURL).Scan(&shortURL)
	if err != nil {
		return ""
	}

	return
}

func (postgresqlDB *PostgresqlDB) GetOriginal(shortURL string) (longURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := postgresqlDB.db.QueryRowContext(ctx, readOriginalURLQuery, shortURL).Scan(&longURL)
	if err != nil {
		return ""
	}

	return
}

func (postgresqlDB *PostgresqlDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := postgresqlDB.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (postgresqlDB *PostgresqlDB) Close() {
	postgresqlDB.db.Close()
}
