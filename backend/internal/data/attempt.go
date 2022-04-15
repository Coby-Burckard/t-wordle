package data

import (
	"database/sql"
	"errors"
	"time"
)

type Attempt struct {
	ID int64
	WordId int64
	// UserId int64 
	IsOpen bool
	IsSolved bool
	CreatedAt time.Time
}

type AttemptModel struct {
	DB *sql.DB
}

// Insert an Attempt into the DB
func (a AttemptModel) Insert(attempt *Attempt) error {
	query := `
		INSERT INTO attempts (word_id)
		VALUES ($1)
		RETURNING id, created_at, is_solved, is_open`

	return a.DB.QueryRow(query, attempt.WordId).Scan(&attempt.ID, &attempt.CreatedAt, &attempt.IsSolved, &attempt.IsOpen)
} 

// Update an Attempt
// Only IsSolved and IsOpen can be modified with this func
func (a AttemptModel) Update(attempt *Attempt) error {
	query := `
		UPDATE attempts
		SET is_solved = $1, is_open = $2
		WHERE id = $3
	`

	return a.DB.QueryRow(query, attempt.IsSolved, attempt.IsOpen, attempt.ID).Scan()
}

// Get an Attempt from the DB
func (a AttemptModel) GetFromWordId(wordId int64) (*Attempt, error) {
	query := `
		SELECT id, word_id, is_open, created_at
		FROM attempts
		WHERE word_id = $1
	`

	var attempt Attempt

	err := a.DB.QueryRow(query, wordId).Scan(
		&attempt.ID,
		&attempt.WordId,
		&attempt.IsOpen,
		&attempt.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows): 				
			return nil, ErrRecordNotFound
		default: 
			return nil, err
		}
	}

	return &attempt, nil
}

// Get an Attempt from the DB if exists, else Create
func (a AttemptModel) GetFromWordIdOrCreate(wordId int64) (*Attempt, error) {
	attempt, err := a.GetFromWordId(wordId)

	// if the attempt is not found, create it
	if err != nil && errors.Is(err, ErrRecordNotFound) {
		attempt = &Attempt{WordId: wordId}
		err = a.Insert(attempt)
	}

	// return any other error or an error resulting from the insert
	if err != nil {
		return nil, err
	}

	return attempt, nil
}

