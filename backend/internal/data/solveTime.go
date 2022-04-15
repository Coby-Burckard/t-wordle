package data

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type SolveTime int

// takes the average time (in ms) and converts it to a
// human readible format _H _M _S
func (t SolveTime) MarshalJSON() ([]byte, error) {
	jsonValue := ""
	remainingSeconds := math.Floor(float64(t/1000))

	hours := math.Floor(float64(remainingSeconds / 60 / 60))
	remainingSeconds -= hours * 60 * 60
	if hours > 0 {
		jsonValue = fmt.Sprintf("%d hours ", int(hours))
	}

	minutes := math.Floor(float64(remainingSeconds / 60))
	remainingSeconds -= minutes * 60
	if minutes > 0 {
		jsonValue += fmt.Sprintf("%d minutes ", int(minutes))
	}

	if remainingSeconds > 0 {
		jsonValue += fmt.Sprintf("%d seconds", int(remainingSeconds))
	}

	return []byte(strconv.Quote(strings.TrimSpace(jsonValue))), nil
}