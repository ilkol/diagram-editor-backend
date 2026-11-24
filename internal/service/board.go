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
	UpdateBoard(ctx context.Context, id, title string, data json.RawMessage) error
}

type BoardService struct {
	storage BoardStorage
}

func NewBoardService(storage BoardStorage) *BoardService {
	return &BoardService{storage: storage}
}

func checkAndNormalizeBoardTitle(title *string) {
	if strings.TrimSpace(*title) == "" {
		*title = "Untitled"
	}
}

func (s *BoardService) Create(ctx context.Context, title string) (*model.Board, error) {
	checkAndNormalizeBoardTitle(&title)

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

func (s *BoardService) Update(ctx context.Context, id, title string, data json.RawMessage) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return model.ErrInvalidBoardID
	}
	checkAndNormalizeBoardTitle(&title)
	if len(data) == 0 {
		data = json.RawMessage("{}")
	}

	return s.storage.UpdateBoard(ctx, id, title, data)
}
