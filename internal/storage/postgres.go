package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ilkol/diagram-editor-backend/internal/model"
	"github.com/jackc/pgx/v5"
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

func (s *Storage) UpdateBoard(ctx context.Context, id, title string, data json.RawMessage) error {
	cmdTag, err := s.pool.Exec(ctx, "UPDATE boards SET title=$1, data=$2, updated_at=$3 WHERE id=$3",
		title, data, time.Now().UTC(), id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return model.ErrBoardNotFound
	}
	return nil
}

func (s *Storage) GetBoard(ctx context.Context, id string) (*model.Board, error) {
	var board model.Board
	err := s.pool.QueryRow(
		ctx,
		"SELECT id, title, data, created_at, updated_at FROM boards WHERE id = $1",
		id,
	).Scan(&board.ID, &board.Title, &board.Data, &board.CreatedAt, &board.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrBoardNotFound
		}
		return nil, err
	}

	return &board, nil
}
