package main

import (
	"fmt"
	"os"
	"trivia/libraries/trivia"
)

func main() {
	err := trivia.Trivia()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
