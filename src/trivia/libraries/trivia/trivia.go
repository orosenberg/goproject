package trivia

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
	"trivia/libraries/prize"
	"trivia/libraries/problems"
)

const limitSeconds = 60
const reward = 10
const penalty = 5

// Trivia game
func Trivia() {
	timeLimit := flag.Int("limit", limitSeconds, "the time limit for the quiz in seconds")
	flag.Parse()

	problems, err := problems.GetProblems()
	if err != nil {
		exit(err)
	}

	// ENTER to start
	fmt.Printf("\n***************** 1980s TV TRIVIA QUIZ *****************\nYou win $%d for every correct answer!\nYou loose $%d for every wrong answer (no worries, you won't go negative, it's all about winning! :D)\nTime limit is %d seconds.\n\nREADY SET GO!(press ENTER to start)\n\n", reward, penalty, limitSeconds)
	reader := bufio.NewReader(os.Stdin)
	start, _ := reader.ReadString('\n')
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	timer.Stop()

	if start[0] == 10 {
		timer.Reset(time.Duration(*timeLimit) * time.Second)
	}
	correct := 0
	readQuestions := 0

	prize := &prize.Prize{0}

problemloop:
	for _, p := range problems {
		readQuestions++

		fmt.Printf("%s\n\n", p.Question)
		fmt.Printf(" 1. %s\n", p.Option1)
		fmt.Printf(" 2. %s\n", p.Option2)
		fmt.Printf(" 3. %s\n\n", p.Option3)
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
			if answer == p.AnswerOptionNumber {
				fmt.Print("CORRECT\n")
				prize.Calculate(reward)
				correct++
				announcement(prize.Amount)
			} else {
				fmt.Printf("INCORRECT => %s. %s\n", p.AnswerOptionNumber, p.AnswerText)
				prize.Calculate(-penalty)
				announcement(prize.Amount)
			}
		}
	}

	fmt.Printf("\nTIME'S UP!\nYou scored %d out of %d. Won $%d.\n", correct, readQuestions, prize.Amount)

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