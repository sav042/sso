package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sav042/sso/internal/domain/models"
	"github.com/sav042/sso/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create conn to db: %s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %s", err.Error())
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uid string, err error) {
	const op = "storage.sqlite.SaveUser"
	userID := uuid.NewString()

	stmt, err := s.db.Prepare("insert into users(id, email, pass_hash) values(?,?,?)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, userID, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (s *Storage) User(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("select id, email, pass_hash from users where email = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var result models.User
	err = stmt.QueryRowContext(ctx, email).Scan(&result.ID, &result.Email, &result.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &result, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID string) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("select is_admin from users where id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var isAdmin bool
	err = stmt.QueryRowContext(ctx, userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID string) (*models.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare("select id, name, secret from apps where id = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var result models.App
	err = stmt.QueryRowContext(ctx, appID).Scan(&result.ID, &result.Name, &result.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &result, nil
}
