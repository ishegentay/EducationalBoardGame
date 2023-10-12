package data

import (
	"EBG.IssataySheg.net/internal/validator"
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

func ValidateMovie(v *validator.Validator, game *Game) {
	v.Check(game.Title != "", "title", "must be provided")
	v.Check(len(game.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(game.Score != 0, "score", "must be provided")
	v.Check(game.Score > 0, "score", "must be a positive integer")
	v.Check(game.Games != nil, "games", "must be provided")
	v.Check(len(game.Games) >= 1, "games", "must contain at least 1 genre")
	v.Check(len(game.Games) <= 5, "games", "must not contain more than 5 genres")
	v.Check(validator.Unique(game.Games), "games", "must not contain duplicate values")
}
