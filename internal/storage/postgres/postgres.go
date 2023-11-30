package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"new-websocket-chat/internal/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(user string, password string, dbname string, hostname string, port int) (*Storage, error) {
	const op = "storage.postgres.New"

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, hostname, port, dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt1, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users(
	    id SERIAL PRIMARY KEY,
	    username CHARACTER VARYING(30) NOT NULL UNIQUE CHECK(username !=''),
	    email CHARACTER VARYING(30) NOT NULL UNIQUE CHECK(email !=''),
		password CHARACTER VARYING(100) NOT NULL);
	`) // migrations ??? not necessary to create table, think about putting in another file and mark in main that you have migrations first
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt1.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt2, err := db.Prepare(`
		CREATE INDEX IF NOT EXISTS idx_username ON users(username);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt2.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(username string, email string, password string) (int64, error) {
	const op = "storage.postgres.SaveUser"

	var lastInsertId int64 = 0
	stmt, err := s.db.Prepare(`INSERT INTO users(username, email, password) VALUES($1, $2, $3) RETURNING id`) // Is it okay? What if table have different structure? How should it work?
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	err = stmt.QueryRow(username, email, password).Scan(&lastInsertId)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == "23505" { // 23505 unique constraint error code - if user already exists
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return lastInsertId, nil
}

func (s *Storage) GetUserEmail(username string) (*string, error) {
	const op = "storage.postgres.GetUserEmail"

	stmt, err := s.db.Prepare(`SELECT email FROM users WHERE username=$1`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resEmail *string
	err = stmt.QueryRow(username).Scan(&resEmail)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrEmailNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resEmail, nil
}

func (s *Storage) GetUsername(email string) (string, error) {
	const op = "storage.postgres.GetUsername"

	stmt, err := s.db.Prepare(`SELECT username FROM users WHERE email=$1`)
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	resUsername := ""
	err = stmt.QueryRow(email).Scan(&resUsername)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrUsernameNotFound
	}

	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resUsername, nil
}

// Returns <nil> if user deleted successfully
func (s *Storage) DeleteUser(username string, email string) error {
	const op = "storage.postgres.DeleteUser"

	stmt, err := s.db.Prepare(`DELETE FROM users WHERE username=$1 AND email=$2`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	result, err := stmt.Exec(username, email)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err) // TODO: why error? level=ERROR msg="failed to detele a user" error="storage.postgres.DeleteUser: execute statement: sql: no rows in result set"
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: gettins rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: no rows affected: %w", op, err)
	}

	return nil
}
