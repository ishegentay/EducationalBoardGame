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

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type GameModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m GameModel) Insert(game *Game) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
		INSERT INTO games (title, score, games)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{game.Title, game.Score, pq.Array(game.Games)}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryRowContext() and pass the context as the first argument.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&game.ID, &game.CreatedAt, &game.Version)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m GameModel) Get(id int64) (*Game, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the movie data.
	query := `
		SELECT  , id, created_at, title, score, games, version
		FROM games
		WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var game Game
	// Use the context.WithTimeout() function to create a context.Context which carries a
	// 3-second timeout deadline. Note that we're using the empty context.Background()
	// as the 'parent' context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&game.ID,
		&game.CreatedAt,
		&game.Title,
		&game.Score,
		pq.Array(&game.Games),
		&game.Version,
	)
	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &game, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m GameModel) Update(game *Game) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
		UPDATE games
		SET title = $1, score = $2, games = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		game.Title,
		game.Score,
		pq.Array(game.Games),
		game.ID,
		game.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryRowContext() and pass the context as the first argument.
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

// Add a placeholder method for deleting a specific record from the movies table.
func (m GameModel) Delete(id int64) error { // Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM games
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// Create a new GetAll() method which returns a slice of movies. Although we're not
// using them right now, we've set this up to accept the various filter parameters as
// arguments.
func (m GameModel) GetAll(title string, Games []string, filters Filters) ([]*Game, Metadata, error) {
	// Construct the SQL query to retrieve all movie records.
	query := fmt.Sprintf(`
			SELECT count(*) OVER(), id, created_at, title, score, games, version
			FROM games
			WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
			AND (games @> $2 OR $2 = '{}')
			ORDER BY %s %s, id ASC
			LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{title, pq.Array(Games), filters.limit(), filters.offset()}
	// And then pass the args slice to QueryContext() as a variadic parameter.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	defer rows.Close()
	// Initialize an empty slice to hold the movie data.
	totalRecords := 0
	games := []*Game{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var game Game
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
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
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Include the metadata struct when returning.
	return games, metadata, nil

}
