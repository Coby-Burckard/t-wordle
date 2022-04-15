package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Attempts AttemptModel
	Guesses GuessModel
	Words WordModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Attempts: AttemptModel{DB: db},
		Guesses: GuessModel{DB: db},
		Words: WordModel{DB: db},
	}
}
