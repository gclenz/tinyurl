package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func GetDatabaseConnection() *sql.DB {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("Unable to connect to databse: %v", err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		slog.Error("Unable to ping databse: %v", err)
		os.Exit(1)
	}

	return db
}

var (
	InsertUrlQuery      = "insert into urls (id, url, created_at, updated_at) values ($1, $2, $3, $4);"
	SelectUrlByIDQuery  = "select id, url, created_at, updated_at from urls where id = $1;"
	UpdateUserCodeQuery = "update urls set url = $1 where id = $2;"
)

type UrlRepository struct {
	db *sql.DB
}

func NewUrlRepository(db *sql.DB) *UrlRepository {
	return &UrlRepository{
		db: db,
	}
}

func (ur *UrlRepository) Create(urlData *UrlData, ctx context.Context) error {
	_, err := ur.db.ExecContext(
		ctx,
		InsertUrlQuery,
		&urlData.ID,
		&urlData.Url,
		&urlData.CreatedAt,
		&urlData.UpdatedAt,
	)

	if err != nil {
		slog.Error("Repository(Create) error:", err)
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicatedKey
		}
		return err
	}

	return nil
}

func (ur *UrlRepository) FindByID(urlID string, ctx context.Context) (*UrlData, error) {
	row := ur.db.QueryRowContext(
		ctx,
		SelectUrlByIDQuery,
		urlID,
	)
	err := row.Err()
	if err != nil {
		slog.Error("Repository(FindByID) error:", err)
		return nil, err
	}

	var urlData UrlData
	err = row.Scan(
		&urlData.ID,
		&urlData.Url,
		&urlData.CreatedAt,
		&urlData.UpdatedAt,
	)
	if err != nil {
		slog.Error("Repository(FindByID) error:", err)
		return nil, err
	}

	return &urlData, nil
}

var (
	ErrDuplicatedKey = errors.New(
		"there id is already in use",
	)
)
