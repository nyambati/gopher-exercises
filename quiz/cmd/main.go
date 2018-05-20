package main

import (
	"flag"

	"github.com/nyambati/go-exercises/quiz"
)

func main() {
	filename := flag.String("csv", "problem.csv", "A csv file in format of question answer")
	limit := flag.Int("limit", 30, "The Limit for taking the quiz")
	flag.Parse()
	quiz := quiz.Quiz{}

	quiz.ReadProblems(*filename)
	quiz.Start(*limit)

}
