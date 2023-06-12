package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var errWrongAnswer = errors.New("wrong answer")

type score struct {
	score int
}

func (s *score) add() {
	s.score++
}

func main() {
	var limit int
	flag.IntVar(&limit, "limit", 30, "limit in seconds to answer the quiz")
	flag.Parse()

	file, err := os.Open("1_quiz_game/problems.csv")
	if err != nil {
		log.Fatal(fmt.Errorf("erro to open file: %w", err))
	}

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(fmt.Errorf("error to read file records: %w", err))
	}

	questions := len(records)

	s := &score{}

	scoreChan := make(chan int)

	questionChan := make(chan int)
	defer close(questionChan)

	timerChan := make(chan int)

	go startScore(scoreChan, timerChan, s)
	go startQuestion(records, scoreChan, questionChan)
	go startTimer(timerChan, questionChan, limit)

	<-timerChan

	fmt.Println("----- SCORE -----")
	fmt.Printf("%d/%d\n", s.score, questions)
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

func startQuestion(records [][]string, scoreChan chan int, questionChan chan int) {
	// wait for timer to start
	<-questionChan
	for _, record := range records {
		err := prompt(record)
		if err == nil {
			scoreChan <- 1
		}
	}
	close(scoreChan)
}

func startScore(scoreChan chan int, timerChan chan int, score *score) {
	for range scoreChan {
		score.add()
	}
	close(timerChan)
}

func startTimer(timerChan chan int, questionChan chan int, secondsLimit int) {
	fmt.Printf("You have %d seconds to answer all questions. Let's play...", secondsLimit)
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(fmt.Errorf("error to read answer: %w", err))
	}

	// start questions
	questionChan <- 1

	<-time.NewTimer(time.Duration(secondsLimit) * time.Second).C
	fmt.Println("\n----- TIME IS OVER -----")
	close(timerChan)
}
