package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ilkol/diagram-editor-backend/internal/model"
)

type BoardService interface {
	Create(ctx context.Context, title string) (*model.Board, error)
}

type BoardHandler struct {
	service BoardService
}

func NewBoardHandler(service BoardService) *BoardHandler {
	return &BoardHandler{service}
}

func (h *BoardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateBoardRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req.Title = ""
	}
	board, err := h.service.Create(r.Context(), req.Title)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(board); err != nil {
		log.Printf("failed to put model data to response: %v", err)
	}
}
