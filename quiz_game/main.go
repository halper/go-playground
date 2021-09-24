package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type QuizItem struct {
	Question string
	Answer   string
}

func NewQuizItem(question, answer string) *QuizItem {
	return &QuizItem{Question: question, Answer: answer}
}

func GetQuizItems(csvFile *os.File) []*QuizItem {
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		os.Exit(-1)
	}
	var items []*QuizItem
	for _, line := range csvLines {
		items = append(items, NewQuizItem(line[0], line[1]))
	}
	return items
}

func Play(score chan int, quizItems []*QuizItem, nCorrect *int) {
	var input string
	for i := 0; i < len(quizItems); i++ {
		q := *quizItems[i]
		fmt.Print(q.Question + " ")
		fmt.Scan(&input)
		if strings.TrimSpace(input) == q.Answer {
			*nCorrect++
		}
		fmt.Println("")
	}
	score <- *nCorrect
}

func main() {
	var f string
	var t int
	currentWorkingD, err := os.Getwd()
	if err != nil {
		fmt.Println("Something went wrong!")
		os.Exit(-1)
	}
	osSep := string(os.PathSeparator)
	fullPath := currentWorkingD + osSep + "problems.csv"
	flag.StringVar(&f, "f", fullPath, "The first number! Default is 1")
	flag.IntVar(&t, "t", 30, "The time in seconds! Default is 30")

	csvFile, err := os.Open(f)
	if err != nil {
		fmt.Println("Can't open the file! " + f)
		os.Exit(-1)
	}
	defer csvFile.Close()

	fmt.Print("Ready to go?")

	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	_, err = reader.ReadString('\n')
	if err != nil {
		os.Exit(-1)
	}
	quizItems := GetQuizItems(csvFile)

	then := time.Now()
	myTimer := time.NewTimer(time.Second * time.Duration(t))

	score := make(chan int)
	defer close(score)
	var nCorrect int
	go Play(score, quizItems, &nCorrect)

	select {
	case <-myTimer.C:
		fmt.Println("\nTimeout!")
	case <-score:
		fmt.Println("Quiz done!")
	}
	n := time.Now()
	fmt.Printf("You answered %d questions correctly out of %d in %d seconds.\n", nCorrect, len(quizItems), int(n.Sub(then).Seconds()))
}
