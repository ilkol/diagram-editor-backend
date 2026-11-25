package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/ilkol/diagram-editor-backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockBoardService struct {
	createFn func(ctx context.Context, title string) (*model.Board, error)
	updateFn func(ctx context.Context, id, title string, data json.RawMessage) error
	getFn    func(ctx context.Context, id string) (*model.Board, error)
}

func (s mockBoardService) Create(ctx context.Context, title string) (*model.Board, error) {
	if s.createFn != nil {
		return s.createFn(ctx, title)
	}
	return nil, nil
}

func (s mockBoardService) Update(ctx context.Context, id, title string, data json.RawMessage) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, title, data)
	}
	return nil
}

func (s mockBoardService) GetByID(ctx context.Context, id string) (*model.Board, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return nil, nil
}

func newRequest(method, body, id string) *http.Request {
	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, "/boards", http.NoBody)
	} else {
		r = httptest.NewRequest(method, "/boards", strings.NewReader(body))
	}

	if id != "" {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", id)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	}
	return r
}

func TestBoardHandler_Create(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		serviceErr  error
		wantStatus  int
		wantTitle   string
		wantSvcCall bool
	}{
		{"valid title", `{"title":"My diagram"}`, nil, http.StatusCreated, "My diagram", true},
		{"empty body", "", nil, http.StatusCreated, "", true},
		{"invalid json", `{"title":`, nil, http.StatusBadRequest, "", false},
		{"service error", `{"title":"My diagram"}`, errors.New("db error"), http.StatusInternalServerError, "My diagram", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotTitle string
			called := false
			mock := mockBoardService{createFn: func(ctx context.Context, title string) (*model.Board, error) {
				called = true
				gotTitle = title
				if tc.serviceErr != nil {
					return nil, tc.serviceErr
				}
				return &model.Board{
					ID:    "board-1",
					Title: title,
					Data:  json.RawMessage(`{}`),
				}, nil
			}}
			h := NewBoardHandler(mock)
			rec := httptest.NewRecorder()

			h.Create(rec, newRequest(http.MethodPost, tc.body, ""))

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, tc.wantSvcCall, called)
			if !tc.wantSvcCall {
				return
			}
			assert.Equal(t, tc.wantTitle, gotTitle)

			if tc.wantStatus == http.StatusCreated {
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
				var board model.Board
				require.NoError(t, json.NewDecoder(rec.Body).Decode(&board))
				assert.Equal(t, "board-1", board.ID)
			}
		})
	}
}

func TestBoardHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		body       string
		serviceErr error
		wantStatus int
		wantTitle  string
		wantData   string
	}{
		{"success", "board-1", `{"title":"Renamed","data":{"nodes":[]}}`, nil, http.StatusOK, "Renamed", `{"nodes":[]}`},
		{"invalid body", "board-1", `{"title":`, nil, http.StatusBadRequest, "", ""},
		{"board not found", "board-1", `{"title":"Renamed"}`, model.ErrBoardNotFound, http.StatusNotFound, "", ""},
		{"invalid board id", "not-a-uuid", `{"title":"Renamed"}`, model.ErrInvalidBoardID, http.StatusBadRequest, "", ""},
		{"service error", "board-1", `{"title":"Renamed"}`, errors.New("db error"), http.StatusInternalServerError, "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotID, gotTitle string
			var gotData json.RawMessage
			mock := mockBoardService{updateFn: func(ctx context.Context, id, title string, data json.RawMessage) error {
				gotID, gotTitle, gotData = id, title, data
				return tc.serviceErr
			}}
			h := NewBoardHandler(mock)
			rec := httptest.NewRecorder()

			h.Update(rec, newRequest(http.MethodPut, tc.body, tc.id))

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantStatus == http.StatusOK {
				assert.Equal(t, tc.id, gotID)
				assert.Equal(t, tc.wantTitle, gotTitle)
				assert.JSONEq(t, tc.wantData, string(gotData))
			}
		})
	}
}

func TestBoardHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		board      *model.Board
		serviceErr error
		wantStatus int
	}{
		{"success", "board-1", &model.Board{ID: "board-1", Title: "Test", Data: json.RawMessage(`{"a":1}`)}, nil, http.StatusOK},
		{"board not found", "board-1", nil, model.ErrBoardNotFound, http.StatusNotFound},
		{"invalid board id", "not-a-uuid", nil, model.ErrInvalidBoardID, http.StatusBadRequest},
		{"service error", "board-1", nil, errors.New("db error"), http.StatusInternalServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotID string
			mock := mockBoardService{getFn: func(ctx context.Context, id string) (*model.Board, error) {
				gotID = id
				return tc.board, tc.serviceErr
			}}
			h := NewBoardHandler(mock)
			rec := httptest.NewRecorder()

			h.GetByID(rec, newRequest(http.MethodGet, "", tc.id))

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, tc.id, gotID)

			if tc.wantStatus != http.StatusOK {
				return
			}
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			var board model.Board
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&board))
			assert.Equal(t, tc.board.ID, board.ID)
			assert.Equal(t, tc.board.Title, board.Title)
			assert.JSONEq(t, string(tc.board.Data), string(board.Data))
		})
	}
}
