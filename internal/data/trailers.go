package data

import (
	"database/sql"
)

type Trailer struct {
	ID           int64  `json:"id"`
	Trailer_name string `json:"trailer_name"`
	Duration     int64  `json:"duration"`
	Premier_date string `json:"premier_date"`
}

type TrailerModel struct {
	DB *sql.DB
}

func (t TrailerModel) Insert(trailer *Trailer) error {
	query := `INSERT INTO trailers (trailer_name, duration, premier_date)
			  VALUES ($1, $2, $3)
			  RETURNING id`

	args := []any{trailer.Trailer_name, trailer.Duration, trailer.Premier_date}

	return t.DB.QueryRow(query, args...).Scan(&trailer.ID)
}
