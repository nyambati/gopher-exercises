package quiz

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// Problem struct
type Problem struct {
	Question string
	Answer   string
}

// Quiz struct Holds methods and data required in the quiz program
type Quiz struct {
	Problems []Problem
}

// ReadProblems : Reads problems from a csv file and parses them
func (q *Quiz) ReadProblems(path string) {
	file, err := os.Open(path)
	if err != nil {
		q.Exit(fmt.Sprintf("Failed to open the CSV file : %s \n", path))
	}
	defer file.Close()

	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		q.Exit("Failed to parse the provided csv file")
	}

	q.Problems = q.ParseProblemToMap(lines)
}

// ParseProblemToMap Parses the read data into a map data type
func (q Quiz) ParseProblemToMap(data [][]string) []Problem {
	problems := make([]Problem, len(data))
	for index, line := range data {
		problems[index] = Problem{
			Question: strings.ToLower(line[0]),
			Answer:   strings.ToLower(strings.TrimSpace(line[1])),
		}
	}
	return problems
}

// Ask asks questions from the generated map
func (q Quiz) Ask(channel chan<- string, qNumber int, question string) {
	fmt.Printf("Problem #%d: %s = ", qNumber+1, question)
	var answer string
	fmt.Scanf("%s\n", &answer)
	channel <- answer
}

// Exit exits quiz when error occurs
func (q Quiz) Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

// Start starts the quiz game
func (q *Quiz) Start(limit int) {
	correct := 0
	total := len(q.Problems)

	timer := time.NewTimer(time.Duration(limit) * time.Second)
	answerChan := make(chan string)

	for i, p := range q.Problems {
		go q.Ask(answerChan, i, p.Question)
		select {
		case <-timer.C:
			q.Score(correct, total)
			return
		case ans := <-answerChan:
			if ans == p.Answer {
				correct++
			}
		}
	}

	q.Score(correct, total)
}

// Score prints the score when time is up or done with the quiz
func (q Quiz) Score(score int, total int) {
	fmt.Println("\n=======================================")
	fmt.Printf("You scored %d out of %d.\n", score, total)
	fmt.Println("=======================================")
}
