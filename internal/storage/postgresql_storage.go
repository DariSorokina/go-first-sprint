package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/DariSorokina/go-first-sprint.git/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	createSchemaQuery = `CREATE SCHEMA IF NOT EXISTS content;`
	createTableQuery  = `CREATE TABLE IF NOT EXISTS content.urls (
		originalURL TEXT, 
		shortURL TEXT,
		userID INTEGER);`
	createIndexQuery      = `CREATE INDEX IF NOT EXISTS originalURL ON content.urls (originalURL)`
	readShortURLQuery     = `SELECT shortURL FROM content.urls WHERE originalURL = $1;`
	readOriginalURLQuery  = `SELECT originalURL FROM content.urls WHERE shortURL = $1;`
	readURLsByUserIDQuery = `SELECT originalURL, shortURL FROM content.urls WHERE userID = $1;`
	writeURLsQuery        = `INSERT INTO content.urls (originalURL, shortURL, userID) VALUES ($1, $2, $3);`
)

type PostgresqlDB struct {
	db *sql.DB
}

func NewPostgresqlDB(cofigBDString string) *PostgresqlDB {
	db, err := sql.Open("pgx", cofigBDString)
	if err != nil {
		log.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, createSchemaQuery)
	if err != nil {
		log.Println(err)
	}
	_, err = db.ExecContext(ctx, createTableQuery)
	if err != nil {
		log.Println(err)
	}
	_, err = db.ExecContext(ctx, createIndexQuery)
	if err != nil {
		log.Println(err)
	}

	return &PostgresqlDB{db: db}
}

func (postgresqlDB *PostgresqlDB) SetValue(shortURL, longURL string, userID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := postgresqlDB.db.ExecContext(ctx, writeURLsQuery, longURL, shortURL, userID)
	if err != nil {
		log.Println(err)
	}

}

func (postgresqlDB *PostgresqlDB) GetShort(longURL string) (shortURL string, errShortURL error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := postgresqlDB.db.QueryRowContext(ctx, readShortURLQuery, longURL).Scan(&shortURL)
	if err != nil {
		return "", nil
	}

	return shortURL, ErrShortURLAlreadyExist
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

func (postgresqlDB *PostgresqlDB) GetURLsByUserID(userID int) (urls []models.URLPair) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := postgresqlDB.db.QueryContext(ctx, readURLsByUserIDQuery, userID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var url models.URLPair
		if err := rows.Scan(&url.OriginalURL, &url.ShortenURL); err != nil {
			log.Fatal(err)
		}
		urls = append(urls, url)
	}

	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
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
	if postgresqlDB.db != nil {
		postgresqlDB.db.Close()
	}
}
