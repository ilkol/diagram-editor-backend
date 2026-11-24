package storage

import (
	"context"

	"github.com/ilkol/diagram-editor-backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}
	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) RunMigrations(ctx context.Context) error {

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(s.pool)
	defer db.Close()
	return goose.UpContext(ctx, db, "migrations")
}

func (s *Storage) CreateBoard(ctx context.Context, board model.Board) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO boards (id, title, data, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		board.ID, board.Title, board.Data, board.CreatedAt, board.UpdatedAt)

	return err
}
