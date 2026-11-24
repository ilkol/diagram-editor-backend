package model

import (
	"encoding/json"
	"errors"
	"time"
)

type Board struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CreateBoardRequest struct {
	Title string `json:"title"`
}

type UpdateBoardRequest struct {
	Title string          `json:"title"`
	Data  json.RawMessage `json:"data"`
}

var ErrBoardNotFound = errors.New("board not found")
var ErrInvalidBoardID = errors.New("invalid board id")
