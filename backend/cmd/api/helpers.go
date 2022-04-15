package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/coby-burckard/t-wordle/backend/internal/data"
	"github.com/julienschmidt/httprouter"
)

// readNumericParam reads an integer id parameter from the request URL.
// It returns the integer value or any error that occured.
func (app *application) readNumericParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	
	return id, nil
}

type envelope map[string]interface{}

func (app *application) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	js = append(js, '\n')
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Configure Decoder
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body to the destination.
	err := dec.Decode(dst)

	// Triage common errors that could occur while decoding JSON
	if err != nil {
			var syntaxError *json.SyntaxError
			var unmarshalTypeError *json.UnmarshalTypeError
			var invalidUnmarshalError *json.InvalidUnmarshalError

			switch {
			case errors.As(err, &syntaxError):
					return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

			case errors.Is(err, io.ErrUnexpectedEOF):
					return errors.New("body contains badly-formed JSON")

			case errors.As(err, &unmarshalTypeError):
					if unmarshalTypeError.Field != "" {
							return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
					}
					return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

			case errors.Is(err, io.EOF):
					return errors.New("body must not be empty")

			case errors.As(err, &invalidUnmarshalError):
					panic(err)

			default:
					return err
			}
	}

	// check request body to ensure no additional and unexpected information has been passed
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
			return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// Responds to the user with the game state
func (app *application) writeBoardJson(w http.ResponseWriter, word *data.Word, attempt *data.Attempt, guesses []*data.Guess) {
	hints := []data.Hint{}
	for _, guess := range guesses {
		_, hint := word.CheckSubmission(guess.Submission)
		hints = append(hints, hint)
	}
	
	data := map[string]interface{}{
		"hints": &hints,
		"attempt": attempt,
	}

	app.writeJson(w, http.StatusCreated, envelope{"data": data}, nil)
}

// getGameState fetches the common game information from the DB specific to the word id and user id.
func (app *application) getGameState(w http.ResponseWriter, r *http.Request, wordId int64) (error, *data.Word, *data.Attempt, []*data.Guess) {
	// get the word from the database
	word := &data.Word{}
	err := app.models.Words.Get(word)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		}	else {
			app.serverErrorResponse(w, r, err)
		}
		return err, nil, nil, nil 
	}

	// get attempt for wordId from the DB if exists, else create
	attempt, err := app.models.Attempts.GetFromWordIdOrCreate(wordId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return err, nil, nil, nil
	}

	// get guesses associated with the attempt
	guesses, err := app.models.Guesses.GetGuessesFromAttemptId(attempt.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return err, nil, nil, nil
	}

	return nil, word, attempt, guesses
}
