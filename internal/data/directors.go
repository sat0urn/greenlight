package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type Director struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Awards  []string `json:"awards,omitempty"`
}

type DirectorModel struct {
	DB *sql.DB
}

func (d DirectorModel) Insert(director *Director) error {
	query := `INSERT INTO directors (name, surname, awards)
			  VALUES ($1, $2, $3)
			  RETURNING id`

	args := []any{director.Name, director.Surname, pq.Array(director.Awards)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return d.DB.QueryRowContext(ctx, query, args...).Scan(&director.ID)
}

func (d DirectorModel) GetAll(name string, filters Filters) ([]*Director, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, name, surname, awards
		FROM directors
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{name, filters.limit(), filters.offset()}

	rows, err := d.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	directors := []*Director{}

	for rows.Next() {
		var director Director

		err := rows.Scan(
			&totalRecords,
			&director.ID,
			&director.Name,
			&director.Surname,
			pq.Array(&director.Awards),
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		directors = append(directors, &director)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return directors, metadata, nil
}
