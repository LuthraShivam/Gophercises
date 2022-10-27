package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

/*
1. CSV file passed with questions in it - CSV entry can be either 1+1,2 or "what 2+2, sir?",4
2. Parse the file.
3. Ask questions to user one by one.

*/
type question struct {
	question string
	answer   string
}

var score int = 0

func parseCSV(csvFileName string) []question {

	questions := []question{}
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to find current working directory, error: %s\n", err)
		return questions
	}

	csvFilePath := currentWorkingDirectory + "/" + csvFileName
	if _, err := os.Stat(csvFilePath); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("File passed does not exist, error: %s\n", err)
		fmt.Println(questions)
		fmt.Println(len(questions))
		return questions
	}

	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("Error in opening CSV file: %s\n", err)
		return questions
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("Error in opening CSV file: %s\n", err)
		return questions
	}

	for _, record := range csvData {
		quizQuestion := question{
			question: record[0],
			answer:   record[1],
		}
		questions = append(questions, quizQuestion)
	}

	return questions
}

func quizTime(quizQuestions []question, timeLimit int) {

	fmt.Println("####### Shivam's Quiz Game #######")

	timeLimitTimer := time.NewTimer(time.Duration(timeLimit) * time.Second) // start timer here
	/*
		Design:
		1. 2 channels - 1 for timer, 1 for input from user
		2. Whatever returns value first, process that and continue
		3. Use channels to get values from these goroutines.
	*/
	answerCh := make(chan string)
	for _, question := range quizQuestions {

		fmt.Printf("Question: %s. Please provide your answer\n", question.question)
		go func() {
			var answer string
			fmt.Scanln(&answer) // this is a blocking statement
			answerCh <- answer
		}()
		select {
		case <-timeLimitTimer.C:
			fmt.Printf("Time limit has run out. You scored %d\n", score)
			return
		case answer := <-answerCh:
			if answer == question.answer {
				score++
			}
		}
	}
	fmt.Printf("You scored: %d\n", score)
}

func main() {

	csvFileName := flag.String("file", "problems.csv", "csv filename used by program")
	timeLimit := flag.Int("time", 30, "quiz time limit (seconds)")
	//shuffle := flag.Bool("shuffle", false, "shuffle questions in the problems file")
	flag.Parse()

	// go func() {
	// 	<-timeLimitTimer.C
	// 	fmt.Printf("Timer has expired. You scored: %d\n", score)
	// }()
	quizQuestions := parseCSV(*csvFileName)
	// if *shuffle {

	// }
	quizTime(quizQuestions, *timeLimit)
}
