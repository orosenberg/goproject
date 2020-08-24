package trivia

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"trivia/libraries/prize"
	"trivia/libraries/problems"
)

const limitSeconds = 60
const reward = 10
const penalty = 5

// Trivia game
func Trivia() error {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,option,option,option,answer'")
	flag.Parse()
	problems, err := problems.GetProblems(*csvFilename)
	if err != nil {
		return err
	}

	// ENTER to start
	fmt.Printf("\n***************** 1980s TV TRIVIA QUIZ *****************\nYou win $%d for every correct answer!\nYou loose $%d for every wrong answer (no worries, you won't go negative, it's all about winning! :D)\nTime limit is %d seconds.\n\nREADY SET GO!(press ENTER to start)\n\n", reward, penalty, limitSeconds)
	reader := bufio.NewReader(os.Stdin)
	start, _ := reader.ReadString('\n')
	timer := time.NewTimer(time.Duration(limitSeconds) * time.Second)
	timer.Stop()

	if start[0] == 10 {
		timer.Reset(time.Duration(limitSeconds) * time.Second)
	}
	correct := 0
	readQuestions := 0

	prize := &prize.Prize{0}

	for _, p := range problems {
		readQuestions++

		fmt.Printf("%s\n\n", p.Question)
		fmt.Printf(" 1. %s\n", p.Option1)
		fmt.Printf(" 2. %s\n", p.Option2)
		fmt.Printf(" 3. %s\n\n", p.Option3)
		fmt.Print("Your answer: ")

		answerCh := make(chan string)
		go getAnswer(answerCh)

		select {
		case <-timer.C:
			gameOver(correct, readQuestions, prize.Amount)
			return nil
		case answer := <-answerCh:
			answerNumber, err := strconv.ParseInt(strings.TrimSpace(answer), 10, 64)
			for err != nil {
				fmt.Println("(Use option numbers as your answer.)")
				go getAnswer(answerCh)
			}
			if answerNumber == p.AnswerOptionNumber {
				fmt.Print("CORRECT\n")
				prize.Calculate(reward)
				correct++
				announcement(prize.Amount)
			} else {
				fmt.Printf("INCORRECT => %d. %s\n", p.AnswerOptionNumber, p.AnswerText)
				prize.Calculate(-penalty)
				announcement(prize.Amount)
			}
		}
	}
	gameOver(correct, readQuestions, prize.Amount)
	return nil
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func getAnswer(answerCh chan string) {
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')

	for !isNumeric(strings.TrimSpace(answer)) {
		fmt.Print("Select one of the options(enter a number): ")
		answer, _ = reader.ReadString('\n')
	}

	answerCh <- answer
}

func gameOver(correct int, readQuestions int, prize int) {
	fmt.Printf("\nTIME'S UP!\nYou scored %d out of %d. Won $%d.\n", correct, readQuestions, prize)
}

func announcement(ammount int) {
	fmt.Printf("-------------- won $%+v so far --------------\n\n", ammount)

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
