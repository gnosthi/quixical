package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {

	//Define flag for specifying csv file.
	//TODO Extend to allow different file formats
	//TODO Extend to add timer flag
	//TODO Extend to add Randomizer flag
	//TODO Add more quiz files covering various Columns, or extend csv to include column numbers.
	csvFile := flag.String("csv", "problems/problems-all.csv", "a csv file in the format of 'question,answer'")
	questions := flag.Int("n", 0, "number of questions to go through")
	flag.Parse()

	// Read in file and parse file and gather questions
	lines, err := readFile(csvFile)
	if err != nil {
		errorExit(fmt.Sprint(err))
	}

	// Set limit defined by -n flag
	if *questions == 0 {
		*questions = len(lines)
	}

	// Create Quiz
	quiz(lines, *questions)

}

// Read and Parse file
func readFile(file *string) ([]problemSet, error) {
	csvFile, err := os.Open(*file)
	if err != nil {
		return []problemSet{}, err
	}

	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	if err != nil {
		return []problemSet{}, err
	}

	problems := parseLines(lines)
	return problems, nil
}

// Run through quiz
func quiz(problems []problemSet, limit int) {
	correct := 0
	for i, p := range problems {
		if limit < i+1 {
			endGame(correct, limit)
		}
		fmt.Printf("Question: #%d: %s = ", i+1, p.question)
		var answer string
		fmt.Scanf("%s\n", &answer)
		if answer == p.answer {
			fmt.Println("Correct!")
			correct++
		} else {
			fmt.Printf("I'm sorry, that is incorrect.\n")
			fmt.Printf("The correct answer was %s\n", p.answer)
		}
	}
}

// End the game
func endGame(correct, total int) {
	fmt.Printf("You got %d out of %d correct!\n", correct, total)
	os.Exit(0)
}

// Parse csv file lines into question and answer sets.
func parseLines(lines [][]string) []problemSet {
	ret := make([]problemSet, len(lines))
	for i, line := range lines {
		ret[i] = problemSet{
			question: line[0],
			answer:   line[1],
		}
	}
	return ret
}

// Create type for question and answer.
type problemSet struct {
	question string
	answer   string
}

// Create an error.
func errorExit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
