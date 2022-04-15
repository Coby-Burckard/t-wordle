package main

import (
	"net/http"

	"github.com/coby-burckard/t-wordle/backend/internal/data"
	"github.com/coby-burckard/t-wordle/backend/internal/validator"
)

// getWordHandler (GET) is an endpoint handlerFunc which returns the users attempt for a specific word.
// A route of /word/:wordId is expected.
// It writes an updated attempt JSON to the user
func (app *application) expireAttemptHandler(w http.ResponseWriter, r *http.Request) {
	// read the word ID
	wordId, err := app.readNumericParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	// Getting word, attempt, and guesses from the user and word id
	err, word, attempt, guesses := app.getGameState(w, r, wordId)
	if err != nil {
		return
	}

	attempt.IsOpen = false
	err = app.models.Attempts.Update(attempt)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// respond to user with updated guess and attempt
	app.writeBoardJson(w, word, attempt, guesses)
}

// guessWordHandler (POST) is an endpoint handlerFunc which handles the creation of a guess.
// A route of /word/:wordId is expected.
// It writes an updated attempt JSON to the user
func (app *application) submitGuessHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Submission string
		SubmissionTime int64
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// read the word ID
	wordId, err := app.readNumericParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return 
	}

	// Deriving word, attempt, and guesses from the user and word id
	_, word, attempt, guesses := app.getGameState(w, r, wordId)
	if err != nil {
		return
	}
	
	// check if the attempt is open, else return the attempt
	if !attempt.IsOpen {
		app.writeBoardJson(w, word, attempt, guesses)
		return 
	}

	// copy input to new Guess instance
	guess := &data.Guess{
		WordId: wordId,
		AttemptId: attempt.ID,
		Submission: input.Submission,
		SubmissionTime: input.SubmissionTime,
	}

	// validate structure of Guess instance
	v := validator.New()
	if data.ValidateGuessRequest(v, guess); v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return 
	}

	// update game state accordingly
	isSolved, _ := word.CheckSubmission(guess.Submission)
	attempt.IsSolved = isSolved
	if len(guesses) >= 4 {
		attempt.IsOpen = false
	} 
	err = app.models.Guesses.Insert(guess)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return 
	}
	app.models.Attempts.Update(attempt)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return 
	}
	guesses = append(guesses, guess)
	
	// respond to user with updated guess and attempt
	app.writeBoardJson(w, word, attempt, guesses)
}

// expireAttemptHandler expires the user's attempt for the given word.
// It writes an updated attempt JSON to the user
func (app *application) getAttemptHandler(w http.ResponseWriter, r *http.Request) {
	// read the word ID
	wordId, err := app.readNumericParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return 
	}

	// Getting word, attempt, and guesses from the user and word id
	_, word, attempt, guesses := app.getGameState(w, r, wordId)
	if err != nil {
		return
	}

	// respond to user with updated guess and attempt
	app.writeBoardJson(w, word, attempt, guesses)
}
