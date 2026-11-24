package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ilkol/diagram-editor-backend/internal/model"
)

type BoardStorage interface {
	CreateBoard(ctx context.Context, board model.Board) error
}

type BoardService struct {
	storage BoardStorage
}

func NewBoardService(storage BoardStorage) *BoardService {
	return &BoardService{storage: storage}
}

func (s *BoardService) Create(ctx context.Context, title string) (*model.Board, error) {
	if strings.TrimSpace(title) == "" {
		title = "Untitled"
	}

	board := model.Board{
		ID:        uuid.NewString(),
		Title:     title,
		Data:      json.RawMessage("{}"),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC()}
	if err := s.storage.CreateBoard(ctx, board); err != nil {
		return nil, err
	}
	return &board, nil
}
