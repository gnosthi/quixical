package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
    "math/rand"
)

var correct int

func main() {

	//Define flag for specifying csv file.
	//TODO Extend to allow different file formats - on hold.
	//TODO Extend to add Randomizer flag
	//TODO Add more quiz files covering various Columns, or extend csv to include column numbers.
	csvFile := flag.String("f", "problems/problems-all.csv", "a csv file in the format of 'question,answer'")
	questions := flag.Int("n", 0, "number of questions to go through")
	timeLimit := flag.Int("t", 0, "Use a timer for the quiz")
    doRandom := flag.Bool("r", true, "Randomize the questions. Default is true")
	flag.Parse()

	// Read in file and parse file and gather questions
	lines, err := readFile(csvFile, *doRandom)

    if err != nil {
		errorExit(fmt.Sprint(err))
	}

	// Set limit defined by -n flag
	if *questions == 0 {
		*questions = len(lines)
	}

	// Create Quiz
	quiz(lines, *questions, timeLimit)

}

func shuffle(lines []problemSet) []problemSet {
    r := rand.New(rand.NewSource(time.Now().Unix()))
    ret := make([]problemSet, len(lines))
    n := len(lines)
    for i := 0; i < n; i++ {
        randIndex := r.Intn(len(lines))
        ret[i] = lines[randIndex]
        lines = append(lines[:randIndex], lines[randIndex+1:]...)
    }
    return ret
}

// Read and Parse file
func readFile(file *string, doRandom bool) ([]problemSet, error) {
	csvFile, err := os.Open(*file)
	if err != nil {
		return []problemSet{}, err
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	if err != nil {
		return []problemSet{}, err
	}
	problems := parseLines(lines, doRandom)
	return problems, nil
}

// Create the timer for quiz
func createTimer(timeLimit int) *time.Timer {
	if timeLimit == 0 {
		fmt.Printf("Running with no time limit\n")
	} else {
		timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
		return timer
	}
	return nil
}

// Ask questions and return whether the answer was correct or incorrect
func quizQuestion(iteration int, question, answer string) string {

	fmt.Printf("Question: #%d: %s = ", iteration, question)
	// Create scanner for input
	input := bufio.NewScanner(os.Stdin)

	// Get user input.
	input.Scan()
	line := input.Text()
	if strings.EqualFold(line, answer) {
		return "Correct"
	}
	return "Wrong"
	/*OLD METHOD
	fmt.Scanf("%s\n", &input)
	if input == answer {
		return "Correct"
	}
	return "Wrong"
	*/
}

//Print function for wrong answers
func wrongAnswer(answer string) {
	fmt.Printf("I'm sorry, that answer is wrong.\n")
	fmt.Printf("The correct answer was %s\n", answer)
}

//Chech answer status and increment correct
func checkAnswer(status, answer string) {
	if status == "Correct" {
		fmt.Println("Correct")
		correct++
	} else {
		wrongAnswer(answer)
	}
}

// Run through quiz
func quiz(problems []problemSet, limit int, timeLimit *int) {
	timer := createTimer(*timeLimit)
	// Ignore timer if -t 0 is set.
	if *timeLimit == 0 {
		timer = nil
	}
	for i, p := range problems {
		if timer != nil {
			// Setup channel for timed questions.
			answerCh := make(chan string)
			go func() {
				q := quizQuestion(i+1, p.question, p.answer)
				answerCh <- q
			}()

			select {
			case <-timer.C:
				endGame(correct, limit)
			case q := <-answerCh:
				checkAnswer(q, p.answer)
			}
		} else {
			if limit < i+1 {
				endGame(correct, limit)
			}
			q := quizQuestion(i+1, p.question, p.answer)
			checkAnswer(q, p.answer)
		}
	}
	endGame(correct, limit)
}

// End the game
func endGame(correct, total int) {
	fmt.Printf("\nYou got %d out of %d correct!\n", correct, total)
	os.Exit(0)
}

// Parse csv file lines into question and answer sets.
func parseLines(lines [][]string, doRandom bool) []problemSet {
	ret := make([]problemSet, len(lines))
	for i, line := range lines {
		ret[i] = problemSet{
			question: line[0],
			answer:   line[1],
		}
	}
    if doRandom == true {
        ret = shuffle(ret)
        return ret
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
