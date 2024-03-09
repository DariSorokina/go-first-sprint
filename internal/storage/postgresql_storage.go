package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/DariSorokina/go-first-sprint.git/internal/logger"
	"github.com/DariSorokina/go-first-sprint.git/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	createSchemaTableIndexQuery = `CREATE SCHEMA IF NOT EXISTS content;
	CREATE TABLE IF NOT EXISTS content.urls (
		originalURL TEXT, 
		shortURL TEXT,
		userID INTEGER,
		deletedFlag BOOLEAN);
	CREATE INDEX IF NOT EXISTS originalURL ON content.urls (originalURL);`
	readShortURLQuery              = `SELECT shortURL FROM content.urls WHERE originalURL = $1;`
	readOriginalURLQuery           = `SELECT originalURL, deletedFlag FROM content.urls WHERE shortURL = $1;`
	readURLsByUserIDQuery          = `SELECT originalURL, shortURL FROM content.urls WHERE userID = $1;`
	writeURLsQuery                 = `INSERT INTO content.urls (originalURL, shortURL, userID, deletedFlag) VALUES ($1, $2, $3, False);`
	updateDeleteFlagQueryBeginning = `UPDATE content.urls SET deletedFlag = True WHERE shortURL in ('`
	updateDeleteFlagQueryEndinning = `') AND userID = ($1);`
)

var ErrReadOriginalURL = errors.New("can not read url")
var ErrDeletedURL = errors.New("requested url was deleted")

type PostgresqlDB struct {
	db  *sql.DB
	log *logger.Logger
}

func NewPostgresqlDB(cofigBDString string, l *logger.Logger) (*PostgresqlDB, error) {
	db, err := sql.Open("pgx", cofigBDString)
	if err != nil {
		l.CustomLog.Sugar().Errorf("Failed to open a database: %s", err)
		return &PostgresqlDB{db: db, log: l}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, createSchemaTableIndexQuery)
	if err != nil {
		l.CustomLog.Sugar().Errorf("Failed to execute a query createSchemaQuery: %s", err)
		return &PostgresqlDB{db: db, log: l}, err
	}

	return &PostgresqlDB{db: db, log: l}, nil
}

func (postgresqlDB *PostgresqlDB) SetValue(ctx context.Context, shortURL, longURL string, userID int) {
	_, err := postgresqlDB.db.ExecContext(ctx, writeURLsQuery, longURL, shortURL, userID)
	if err != nil {
		postgresqlDB.log.CustomLog.Sugar().Errorf("Failed to execute a query writeURLsQuery: %s", err)
	}

}

func (postgresqlDB *PostgresqlDB) GetShort(ctx context.Context, longURL string) (shortURL string, errShortURL error) {
	err := postgresqlDB.db.QueryRowContext(ctx, readShortURLQuery, longURL).Scan(&shortURL)
	if err != nil {
		return "", nil
	}

	return shortURL, ErrShortURLAlreadyExist
}

func (postgresqlDB *PostgresqlDB) GetOriginal(ctx context.Context, shortURL string) (longURL string, getOriginalErr error) {
	var deletedFlag bool

	err := postgresqlDB.db.QueryRowContext(ctx, readOriginalURLQuery, shortURL).Scan(&longURL, &deletedFlag)
	if err != nil {
		return "", ErrReadOriginalURL
	}

	if deletedFlag {
		return "", ErrDeletedURL
	}

	return longURL, nil
}

func (postgresqlDB *PostgresqlDB) GetURLsByUserID(ctx context.Context, userID int) (urls []models.URLPair) {

	rows, err := postgresqlDB.db.QueryContext(ctx, readURLsByUserIDQuery, userID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var url models.URLPair
		if err := rows.Scan(&url.OriginalURL, &url.ShortenURL); err != nil {
			postgresqlDB.log.CustomLog.Sugar().Errorf("Failed to scan original and shorten urls in GetURLsByUserID method: %s", err)
		}
		urls = append(urls, url)
	}

	rerr := rows.Close()
	if rerr != nil {
		postgresqlDB.log.CustomLog.Sugar().Errorf("Close error in GetURLsByUserID method: %s", rerr)
	}

	if err := rows.Err(); err != nil {
		postgresqlDB.log.CustomLog.Sugar().Errorf("The last error encountered by Rows.Scan in GetURLsByUserID method: %s", err)
		log.Fatal(err)
	}

	return
}

func (postgresqlDB *PostgresqlDB) Ping(ctx context.Context) error {
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

func (postgresqlDB *PostgresqlDB) DeleteURLsWorker(shortURLs []string, userID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	processedShortURLs := strings.Join(shortURLs, "', '")
	updateDeleteFlagQuery := updateDeleteFlagQueryBeginning + processedShortURLs + updateDeleteFlagQueryEndinning

	result, err := postgresqlDB.db.ExecContext(ctx, updateDeleteFlagQuery, userID)
	if err != nil {
		postgresqlDB.log.CustomLog.Sugar().Errorf("Failed to execute a query updateDeleteFlagQuery: %s", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		postgresqlDB.log.CustomLog.Sugar().Errorf("Failed to execute RowsAffected: %s", err)
	}
	if rows != 1 {
		postgresqlDB.log.CustomLog.Sugar().Infof("Affected rows: %d", rows)
	}
}
