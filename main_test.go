package main

import (
	"flag"
	"testing"
    "os"
    "os/exec"
)

var csvFile = flag.String("csv", "test_problems.csv", "csv test")

func TestReadFile(t *testing.T) {
    doRandom := false
	lines, err := readFile(csvFile, doRandom)
	if err != nil {
		t.Error(err)
	}
	switchQAndA(t, lines)
}

func BenchmarkReadFile(t *testing.B) {
    doRandom := false
    lines, err := readFile(csvFile, doRandom)
    if err != nil {
        t.Error(err)
    }
    benchSwitchQAndA(t, lines)
}

func TestParseLines(t *testing.T) {
    doRandom := false
	problems := [][]string{
		{"Question?", "Answer"},
		{"1+2", "3"},
		{"Foo", "Bar"},
		{"Alice", "Bob"},
		{"Mock", "Spock"},
	}
	lines := parseLines(problems, doRandom)
	switchQAndA(t, lines)
}

func BenchmarkParseLines(t *testing.B) {
    doRandom := false
    problems := [][]string{
        {"Question?", "Answer"},
        {"1+2", "3"},
        {"Foo", "Bar"},
        {"Alice", "Bob"},
        {"Mock", "Spock"},
    }
    lines := parseLines(problems, doRandom)
    benchSwitchQAndA(t, lines)
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

func benchSwitchQAndA(t *testing.B, lines []problemSet) {
    for i, p := range lines {
        if len(lines) <= 1{
            t.Errorf("Total number of questions is %d. Expected %d", i, 5)
        }
        switch i {
        case 0:
            benchQAParseError(t, "Question?", p.question)
            benchQAParseError(t, "Answer", p.answer)
        case 1:
            benchQAParseError(t, "1+2", p.question)
            benchQAParseError(t, "3", p.answer)
        case 2:
            benchQAParseError(t, "Foo", p.question)
            benchQAParseError(t, "Bar", p.answer)
        case 3:
            benchQAParseError(t, "Alice", p.question)
            benchQAParseError(t, "Bob", p.answer)
        case 4:
            benchQAParseError(t, "Mock", p.question)
            benchQAParseError(t, "Spock", p.answer)
        }
    }
}

func qaParseError(t *testing.T, expected, recieved string) {
	if expected != recieved {
		t.Errorf("Expected %s, got %s", expected, recieved)
	}
}

func benchQAParseError(t *testing.B, expected, recieved string) {
    if expected != recieved {
        t.Errorf("Expected %s, got %s", expected, recieved)
    }
}
