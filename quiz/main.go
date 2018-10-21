package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Problem struct {
	question string
	answer   int
}

const PROBLEM_BUFFER_COUNT int = 100
const DEFAULT_TIME_LIMIT int = 30

var count, score, faults int

func readProblems(probems chan Problem, filename string, shuffle bool) {
	buffer := make([]Problem, PROBLEM_BUFFER_COUNT)
	data, err := os.Open(filename)

	onError(err)

	defer data.Close()

	scanner := bufio.NewScanner(data)

	fmt.Print(scanner.Scan())

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		question := line[0]

		answer, err := strconv.Atoi(line[1])

		onError(err)

		buffer[count] = Problem{question, answer}
		count++
	}

	if shuffle {
		rand.Seed(time.Now().Unix())
		for i := range buffer[:count] {
			j := rand.Intn(i + 1)
			buffer[i], buffer[j] = buffer[j], buffer[i]
		}
	}

	for _, problem := range buffer[:count] {
		probems <- problem
	}
}

func solveProblem(problems chan Problem) {
	scanner := bufio.NewScanner(os.Stdin)
	for problem := range problems {
		fmt.Printf("%s = ", problem.question)

		scanner.Scan()

		input := strings.Trim(scanner.Text(), " ")

		answer, err := strconv.Atoi(input)

		if err != nil {
			fmt.Printf(" '%s' is not a valid answer \n", input)
			continue
		}

		if answer == problem.answer {
			score++
			fmt.Println("Thats the corrent answer !")
		} else {
			faults++
			fmt.Println("Your answer is wrong!")
		}
	}

	return
}

func startTimer(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
func onError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	filename := flag.String("file", "problem.csv", "Location of the question csv file")
	seconds := flag.Int("seconds", DEFAULT_TIME_LIMIT, "Number of seconds to solve problems")
	shuffle := flag.Bool("shuffle", false, "Shuffle questions")
	debug := flag.Bool("debug", false, "Show debug information")
	flag.Parse()

	if !*debug {
		log.SetOutput(ioutil.Discard)
	}

	log.Println("debug:", *debug)
	log.Println("filename:", *filename)
	log.Println("shuffle:", *shuffle)
	log.Println("timer:", *seconds)

	problems := make(chan Problem, PROBLEM_BUFFER_COUNT)

	go readProblems(problems, *filename, *shuffle)
	fmt.Printf("Press any key to start the quiz!")
	bufio.NewScanner(os.Stdin).Scan()
	go solveProblem(problems)

	startTimer(*seconds)

	fmt.Printf("\nNumber of questions: %d\n", count)
	fmt.Printf("Correct answers: %d\n", score)
	fmt.Printf("Incorrect answers: %d\n", faults)
}
