package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
)

type Word struct {
	ID int64
	Answer string
	SolveCount int
	SolveTime int 
}

type charHint struct {
	char string
	check int 
}

func (c charHint) MarshalJSON() ([]byte, error) {
	jsonValue := map[string]interface{}{}
	jsonValue["char"] = c.char
	jsonValue["check"] = c.check
	return json.Marshal(jsonValue)
}

type Hint [5]charHint

func (w Word) CheckSubmission(submission string) (bool, Hint) {
	isSolved := false
	hint := Hint{}
	answer := w.Answer

	// find exact matches and remove from answer
	for i, character := range submission {
		hint[i].char = string(character)
		if string(answer[i]) == hint[i].char {
			hint[i].check = 2
			rAnswer := []rune(answer)
			rAnswer[i] = '-'
			answer = string(rAnswer)
		} 
	}

	// find char in word that are not exact matches
	for i, _ := range submission {
		if strings.Contains(answer, hint[i].char) {
			hint[i].check = 1
			answer = strings.Replace(answer, hint[i].char, "-", 1)
		} 
	}

	// checksum on hint
	checksum := 0
	for _, charCheck := range hint {
		checksum += charCheck.check
	}
	if (checksum == 10) {
		isSolved = true
	}
	return isSolved, hint
}

type WordModel struct {
	DB *sql.DB
}

func (w WordModel) Get(word *Word) error {
	query := `
		SELECT answer, solve_count, solve_time
		FROM words
		WHERE id = $1`

	row := w.DB.QueryRow(query, word.ID)
	err := row.Scan(&word.Answer, &word.SolveCount, &word.SolveTime)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows): 				
			return ErrRecordNotFound
		default: 
			return err
		}
	}

	return nil
}