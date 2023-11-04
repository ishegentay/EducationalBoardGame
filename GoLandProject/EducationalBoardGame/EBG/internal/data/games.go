package data

import (
	"EBG.IssataySheg.net/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
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

type GameModel struct {
	DB *sql.DB
}

func (m GameModel) Insert(game *Game) error {
	query := `
		INSERT INTO games (title, score, games)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`
	args := []interface{}{game.Title, game.Score, pq.Array(game.Games)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&game.ID, &game.CreatedAt, &game.Version)
}

func (m GameModel) Get(id int64) (*Game, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT  , id, created_at, title, score, games, version
		FROM games
		WHERE id = $1`
	var game Game
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&game.ID,
		&game.CreatedAt,
		&game.Title,
		&game.Score,
		pq.Array(&game.Games),
		&game.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &game, nil
}

func (m GameModel) Update(game *Game) error {
	query := `
		UPDATE games
		SET title = $1, score = $2, games = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version`
	args := []interface{}{
		game.Title,
		game.Score,
		pq.Array(game.Games),
		game.ID,
		game.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&game.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m GameModel) Delete(id int64) error { // Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM games
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m GameModel) GetAll(title string, Games []string, filters Filters) ([]*Game, Metadata, error) {
	query := fmt.Sprintf(`
			SELECT count(*) OVER(), id, created_at, title, score, games, version
			FROM games
			WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
			AND (games @> $2 OR $2 = '{}')
			ORDER BY %s %s, id ASC
			LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{title, pq.Array(Games), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	games := []*Game{}
	for rows.Next() {
		var game Game
		err := rows.Scan(
			&totalRecords,
			&game.ID,
			&game.CreatedAt,
			&game.Title,
			&game.Score,
			pq.Array(&game.Games),
			&game.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		games = append(games, &game)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return games, metadata, nil

}
