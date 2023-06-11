package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var errWrongAnswer = errors.New("wrong answer")

func main() {
	file, err := os.Open("1_quiz_game/problems.csv")
	if err != nil {
		log.Fatal(fmt.Errorf("erro to open file: %w", err))
	}

	reader := csv.NewReader(file)
	questions := 0
	score := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(fmt.Errorf("error to read file line: %w", err))
		}

		questions++
		err = prompt(record)
		if err == nil {
			score++
		}
	}

	fmt.Println("----- SCORE -----")
	fmt.Printf("%d/%d\n", score, questions)
}

func prompt(record []string) error {
	question := strings.TrimSpace(record[0])
	resp := strings.TrimSpace(record[1])

	fmt.Printf("%s = ", question)

	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(fmt.Errorf("error to read answer: %w", err))
	}

	answer = strings.TrimSpace(answer)

	if answer != resp {
		return errWrongAnswer
	}

	return nil
}
