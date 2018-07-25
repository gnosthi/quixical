package main

import (
	"flag"
	"testing"
)

func TestReadFile(t *testing.T) {
	csvFile := flag.String("csv", "test_problems.csv", "csv test")
	lines, err := readFile(csvFile)
	if err != nil {
		t.Error(err)
	}
	switchQAndA(t, lines)
}

func TestParseLines(t *testing.T) {
	problems := [][]string{
		{"Question?", "Answer"},
		{"1+2", "3"},
		{"Foo", "Bar"},
		{"Alice", "Bob"},
		{"Mock", "Spock"},
	}
	lines := parseLines(problems)
	switchQAndA(t, lines)
}

func switchQAndA(t *testing.T, lines []problemSet) {
	for i, p := range lines {
		if len(lines) <= 1 {
			t.Errorf("Total number of questions is %d. Expected %d", i, 5)
		}
		switch i {
		case 0:
			qaParseError(t, "Question?", p.question)
			qaParseError(t, "Answer", p.answer)
		case 1:
			qaParseError(t, "1+2", p.question)
			qaParseError(t, "3", p.answer)
		case 2:
			qaParseError(t, "Foo", p.question)
			qaParseError(t, "Bar", p.answer)
		case 3:
			qaParseError(t, "Alice", p.question)
			qaParseError(t, "Bob", p.answer)
		case 4:
			qaParseError(t, "Mock", p.question)
			qaParseError(t, "Spock", p.answer)
		default:
			break
		}
	}
}

func qaParseError(t *testing.T, expected, recieved string) {
	if expected != recieved {
		t.Errorf("Expected %s, got %s", expected, recieved)
	}
}
