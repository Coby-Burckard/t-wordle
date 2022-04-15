package data

import (
	"database/sql"
	"time"

	"github.com/coby-burckard/t-wordle/backend/internal/validator"
)

type Guess struct {
	ID int64
	WordId int64
	// UserId int64
	AttemptId int64
	Submission string
	CreatedAt time.Time
	SubmissionTime int64
}


func ValidateGuessRequest(v *validator.Validator, guess *Guess) {
	// check that a guess.Submission word exists and is 5 letters
	v.Check(guess.Submission != "", "Guess", "Guess cannot be empty")
	v.Check(len(guess.Submission) != 5, "Guess", "Guess must be 5 letters") 
	
	// check that a required fields were provided
	v.Check(guess.WordId > 0, "Word", "WordId must be provided")
	v.Check(guess.SubmissionTime > 0, "Submission Time", "SubmissionTime must be provided")
}

type GuessModel struct {
	DB *sql.DB	
}

func (g GuessModel) GetGuessesFromAttemptId(attemptId int64) ([]*Guess, error) {
	query := `
		SELECT id, word_id, attempt_id, submission, created_at, submission_time
		FROM guesses
		WHERE attempt_id = $1
	`

	rows, err := g.DB.Query(query, attemptId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load the matched rows into a slice of guesses
	guesses := []*Guess{}
	for rows.Next() {
		guess := &Guess{}
		err = rows.Scan(&guess.ID, &guess.WordId, &guess.AttemptId, &guess.Submission, &guess.CreatedAt, &guess.SubmissionTime)
		if err != nil {
			return nil, err
		}
		guesses = append(guesses, guess)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return guesses, nil
}

func (g GuessModel) Insert(guess *Guess) error {
	query := `
		INSERT INTO guesses (word_id, attempt_id, submission, submission_time)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	return g.DB.QueryRow(query, guess.WordId, guess.AttemptId, guess.Submission, guess.SubmissionTime).Scan(&guess.ID, &guess.CreatedAt)
}
