package data

import (
	"time"
)

type Game struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Games       []string  `json:"games,omitempty"`
	Score       Score     `json:"score,omitempty"`
	Version     int32     `json:"version"`
}
