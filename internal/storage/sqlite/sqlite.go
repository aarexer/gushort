package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"gushort/internal/storage"
	"time"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

const saveUrlSql = `
		insert into
			url_shorten(
				url,
				alias,
				created_at,
				updated_at
			)
		values
			(
				?,
				?,
				?,
				?
			)
	`

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare(saveUrlSql)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	now := time.Now()
	res, err := stmt.Exec(urlToSave, alias, now, now)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUrlAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

const getUrlByAliasSql = `
	select
		url
	from
		url_shorten
	where
		alias = ?
	`

func (s *Storage) GetUrlByAlias(alias string) (string, error) {
	const op = "storage.sqlite.GetUrlByAlias"

	stmt, err := s.db.Prepare(getUrlByAliasSql)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var result string
	err = stmt.QueryRow(alias).Scan(&result)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return result, nil
}
