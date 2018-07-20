package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {

	//Define flag for specifying csv file.
	csvFile := flag.String("csv", "problems/problems-all.csv", "a csv file in the format of 'question,answer'")
	questions := flag.Int("n", 0, "number of questions to go through")
	flag.Parse()

	//Reference csv file by pointer
	file, err := os.Open(*csvFile)
	if err != nil {
		errorExit("Failed to open file: " + *csvFile)
	}

	//Read in csvFile.
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		errorExit("Failed to parse csv content in: " + *csvFile)
	}

	// Load problem set and modify number of questions
	problems := parseLines(lines)
	if *questions == 0 {
		*questions = len(problems)
	}

	correct := 0
	wrong := 0
	// Iterate over problem sets.
	for i, p := range problems {
		if *questions < i+1 {
			os.Exit(0)
		}
		fmt.Printf("Question: #%d: %s = ", i+1, p.question)
		var answer string
		fmt.Scanf("%s\n", &answer)
		if answer == p.answer {
			fmt.Println("Correct!")
			correct++
		} else {
			fmt.Println("Wrong!")
			fmt.Printf("The correct answer was %s\n", p.answer)
			wrong++
		}
	}
	fmt.Printf("You got %d out of %d correct!\n", correct, len(problems))
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
