package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ilkol/diagram-editor-backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	createFn func(ctx context.Context, board model.Board) error
	getFn    func(ctx context.Context, id string) (*model.Board, error)
	updateFn func(ctx context.Context, id, title string, data json.RawMessage) error
}

func (s mockStorage) CreateBoard(ctx context.Context, board model.Board) error {
	if s.createFn != nil {
		return s.createFn(ctx, board)
	}
	return nil
}
func (s mockStorage) UpdateBoard(ctx context.Context, id, title string, data json.RawMessage) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, title, data)
	}
	return nil
}
func (s mockStorage) GetBoard(ctx context.Context, id string) (*model.Board, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return nil, nil
}

func TestBoardService_Create(t *testing.T) {
	tests := []struct {
		name       string
		inputTitle string
		mockErr    error
		wantTitle  string
		wantErr    bool
	}{
		{"valid title", "Some diagram", nil, "Some diagram", false},
		{"empty title", "", nil, "Untitled", false},
		{"space title", "   ", nil, "Untitled", false},
		{"bd error", "", errors.New("db error"), "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockStorage{createFn: func(ctx context.Context, board model.Board) error {
				return tc.mockErr
			}}
			svc := NewBoardService(mock)
			board, err := svc.Create(context.Background(), tc.inputTitle)
			if tc.wantErr == false {
				require.NoError(t, err)
				require.NotNil(t, board)
				assert.Equal(t, tc.wantTitle, board.Title)
				_, parseErr := uuid.Parse(board.ID)
				assert.NoError(t, parseErr)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestBoardService_GetByID(t *testing.T) {
	validID := uuid.NewString()
	tests := []struct {
		name    string
		id      string
		mockErr error
		wantErr bool
		storage *mockStorage
	}{
		{"invalid id", "not-a-uuid", model.ErrInvalidBoardID, true, nil},
		{"board not found", uuid.NewString(), model.ErrBoardNotFound, true, &mockStorage{getFn: func(ctx context.Context, id string) (*model.Board, error) { return nil, model.ErrBoardNotFound }}},
		{"success", validID, nil, false, &mockStorage{getFn: func(ctx context.Context, id string) (*model.Board, error) {
			return &model.Board{ID: validID, Title: "Test"}, nil
		}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewBoardService(tc.storage)
			board, err := svc.GetByID(context.Background(), tc.id)
			if tc.wantErr == false {
				require.NoError(t, err)
				require.NotNil(t, board)
				assert.Equal(t, tc.id, board.ID)
				assert.Equal(t, "Test", board.Title)
				_, parseErr := uuid.Parse(board.ID)
				assert.NoError(t, parseErr)
			} else {
				require.ErrorIs(t, err, tc.mockErr)
			}
		})
	}

}
