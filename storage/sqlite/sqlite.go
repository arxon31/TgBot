package sqlite

import (
	"TgBot/storage"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

//TODO go doc comments for each function

func New(path string) (*Storage, error) {
	//Opening connection to our database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	//Ping our database to ensure that we are connected
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil

}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?,?)`
	_, err := s.db.ExecContext(ctx, q, p.URL, p.Username)
	if err != nil {
		return fmt.Errorf("can't save page to database: %w", err)
	}
	return nil
}

func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page : %w", err)
	}

	return &storage.Page{
		URL:      url,
		Username: userName,
	}, nil
}

func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, p.URL, p.Username); err != nil {
		return fmt.Errorf("can't delete page: %w", err)
	}
	return nil
}

func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	err := s.db.QueryRowContext(ctx, q, p.URL, p.Username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`
	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}
	return nil
}
