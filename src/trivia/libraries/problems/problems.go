package problems

import (
	"encoding/csv"
	"flag"
	"os"
	"strconv"
	"strings"
)

// GetProblems returns trivia questions
func GetProblems() ([]Problem, error) {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,option,option,option,answer'")
	flag.Parse()

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
		solutionOptionNumber, err := strconv.ParseInt(strings.TrimSpace(line[4]), 10, 64)
		if err != nil {
			return nil, err
		}

		ret[i] = Problem{
			Question:           line[0],
			Option1:            line[1],
			Option2:            line[2],
			Option3:            line[3],
			AnswerOptionNumber: strings.TrimSpace(line[4]),
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
	AnswerOptionNumber string
	AnswerText         string
}
