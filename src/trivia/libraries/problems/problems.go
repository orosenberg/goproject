package problems

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GetProblems returns trivia questions
func GetProblems(csvFilename *string) ([]Problem, error) {
	file, err := os.Open(*csvFilename)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return parseLines(lines)
}

func parseLines(lines [][]string) ([]Problem, error) {
	ret := make([]Problem, len(lines))
	for i, line := range lines {
		if len(line) != 5 {
			err := fmt.Errorf("Wrong number of data columns, expected 5, received %d", len(line))
			return nil, err
		}

		solutionOptionNumber, err := strconv.ParseInt(strings.TrimSpace(line[4]), 10, 64)
		if err != nil {
			return nil, err
		}

		if !(solutionOptionNumber > 0 && solutionOptionNumber < 4) {
			err = fmt.Errorf("Solution number should be between 1-3 inclusive, received %d", solutionOptionNumber)
			return nil, err
		}

		ret[i] = Problem{
			Question:           line[0],
			Option1:            line[1],
			Option2:            line[2],
			Option3:            line[3],
			AnswerOptionNumber: solutionOptionNumber,
			AnswerText:         line[solutionOptionNumber],
		}
	}
	return ret, nil
}

// Problem trivia
type Problem struct {
	Question           string
	Option1            string
	Option2            string
	Option3            string
	AnswerOptionNumber int64
	AnswerText         string
}
