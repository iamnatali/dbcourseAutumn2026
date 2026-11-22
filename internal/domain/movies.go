package domain

import (
	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	ReleaseDate Date      `json:"releaseDate" db:"release_date"`
	//moview int64 runtime in minutes
	//todo: replace with time.Duration
	//Runtime time.Duration  `json:"runtime,omitempty" db:"runtime"`
	Runtime  Runtime  `json:"runtime" db:"runtime"`
	Genres   []string `json:"genres" db:"genres"`
	Language string `json:"language" db:"language"`
}
