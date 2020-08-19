package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const limitSeconds = 60
const reward = 10
const penalty = 5

type prize struct {
	amount int
}

func (p *prize) calculate(change int) {
	if p.amount+change < 0 {
		p.amount = 0
	} else {
		p.amount += change
	}
}

func announcement(ammount int) {
	fmt.Printf("\n-------------- won $%+v so far --------------\n\n", ammount)
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,option,option,option,answer'")
	timeLimit := flag.Int("limit", limitSeconds, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	problems := parseLines(lines)

	// ENTER to start
	fmt.Printf("\n***************** 1980s TV TRIVIA QUIZ *****************\nYou win $%d for every correct answer!\nYou loose $%d for every wrong answer (no worries, you won't go negative, it's all about winning! :D)\nTime limit is %d seconds.\n\nREADY SET GO!(press ENTER to start)", reward, penalty, limitSeconds)
	reader := bufio.NewReader(os.Stdin)
	start, _ := reader.ReadString('\n')
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	timer.Stop()

	if start[0] == 10 {
		timer.Reset(time.Duration(*timeLimit) * time.Second)
	}
	correct := 0
	readQuestions := 0

	prize := &prize{0}

problemloop:
	for _, p := range problems {
		readQuestions++
		announcement(prize.amount)

		fmt.Printf("%s\n\n", p.q)
		fmt.Printf(" 1. %s\n", p.o1)
		fmt.Printf(" 2. %s\n", p.o2)
		fmt.Printf(" 3. %s\n\n", p.o3)
		fmt.Print("Your answer: ")

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-answerCh:
			if answer == p.a {
				fmt.Print("*** CORRECT ***\n\n")
				prize.calculate(reward)
				correct++
				fmt.Print("NEXT (press ENTER to continue)")
				reader := bufio.NewReader(os.Stdin)
				_, _ = reader.ReadString('\n')
			} else {
				fmt.Printf("*** INCORRECT *** => %s. %s\n\n", p.a, p.aText)
				prize.calculate(-penalty)
				fmt.Print("NEXT (press ENTER to continue)")
				reader := bufio.NewReader(os.Stdin)
				_, _ = reader.ReadString('\n')
			}
		}
	}

	fmt.Printf("\nTIME'S UP!\nYou scored %d out of %d. Won $%d.\n", correct, readQuestions, prize.amount)
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		solutionOptionNumber, _ := strconv.ParseInt(strings.TrimSpace(line[4]), 10, 64)

		ret[i] = problem{
			q:     line[0],
			o1:    line[1],
			o2:    line[2],
			o3:    line[3],
			a:     strings.TrimSpace(line[4]),
			aText: line[solutionOptionNumber],
		}
	}
	return ret
}

type problem struct {
	q     string
	o1    string
	o2    string
	o3    string
	a     string
	aText string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
