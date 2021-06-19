package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// Get returns one movie and error, if any
func (m *DBModel) Get(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, title, description, year, release_date, rating, runtime, mpaa_rating,
				created_at, updated_at from movies where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	var movie Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Rating,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// get the genres
	query = `select mg.id, mg.movie_id, mg.genre_id, mg.created_at, mg.updated_at,
       				g.id, g.genre_name, g.created_at, g.updated_at
				from movies_genres mg
				left join genres g on g.id = mg.genre_id
				where movie_id = $1`
	rows, _ := m.DB.QueryContext(ctx, query, id)
	defer rows.Close()

	var movieGenres []MovieGenre
	for rows.Next() {
		var mg MovieGenre
		err = rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.CreatedAt,
			&mg.UpdatedAt,
			&mg.Genre.ID,
			&mg.Genre.GenreName,
			&mg.Genre.CreatedAt,
			&mg.Genre.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		movieGenres = append(movieGenres, mg)
	}

	movie.MovieGenre = movieGenres

	return &movie, nil
}

// All returns all movies and error, if any
func (m *DBModel) All(id int) ([]*Movie, error) {
	return nil, nil
}