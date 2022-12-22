package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/mheidinger/advent-of-code-go/cast"
	"github.com/mheidinger/advent-of-code-go/util"
)

//go:embed input.txt
var input string

func init() {
	// do this in init (not main) so test file has same input
	input = strings.TrimRight(input, "\n")
	if len(input) == 0 {
		panic("empty input.txt file")
	}
}

func main() {
	var part int
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.Parse()
	fmt.Println("Running part", part)

	if part == 1 {
		ans := part1(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

type Statement struct {
	Name       string
	Num        int
	Operator   string
	Operand1   string
	Operand2   string
	Resolved   bool
	Statement1 *Statement
	Statement2 *Statement
	Human      bool
}

func Resolve(statementMap map[string]*Statement, name string) *Statement {
	statement, ok := statementMap[name]
	if !ok {
		panic(fmt.Errorf("Unknown name: %s", name))
	}
	if statement.Resolved || statement.Operator == "" || statement.Human {
		return statement
	}
	operand1 := Resolve(statementMap, statement.Operand1)
	operand2 := Resolve(statementMap, statement.Operand2)
	if operand1.Operator == "" && operand2.Operator == "" {
		switch statement.Operator {
		case "*":
			statement.Num = operand1.Num * operand2.Num
		case "+":
			statement.Num = operand1.Num + operand2.Num
		case "/":
			statement.Num = operand1.Num / operand2.Num
		case "-":
			statement.Num = operand1.Num - operand2.Num
		default:
			panic(fmt.Errorf("Unknown operator %s", statement.Operator))
		}
		statement.Operator = ""
	} else {
		statement.Statement1 = operand1
		statement.Statement2 = operand2
	}

	statement.Resolved = true
	return statement
}

func part1(input string) int {
	statements := parseInput(input)

	statementMap := make(map[string]*Statement, len(statements))
	for _, statement := range statements {
		statementMap[statement.Name] = statement
	}

	return Resolve(statementMap, "root").Num
}

func resetMap(statementMap map[string]*Statement) {
	for _, statement := range statementMap {
		statement.Resolved = false
	}
}

func getFixAndSolve(stmt *Statement) (int, *Statement) {
	fixedNum := -1
	var solveStmt *Statement
	if stmt.Statement1.Operator == "" {
		fixedNum = stmt.Statement1.Num
		solveStmt = stmt.Statement2
	} else if stmt.Statement2.Operator == "" {
		fixedNum = stmt.Statement2.Num
		solveStmt = stmt.Statement1
	} else {
		panic(fmt.Errorf("No fixed num in statement: %+v", stmt))
	}

	return fixedNum, solveStmt
}

func SolveToResult(stmt *Statement, want int) int {
	if stmt.Human {
		return want
	}

	solveWant := -1
	var solveStmt *Statement

	if stmt.Statement1.Operator == "" {
		fixedNum := stmt.Statement1.Num
		solveStmt = stmt.Statement2

		switch stmt.Operator {
		case "*":
			solveWant = want / fixedNum
		case "+":
			solveWant = want - fixedNum
		case "/":
			solveWant = fixedNum / want
		case "-":
			solveWant = fixedNum - want
		}
	} else if stmt.Statement2.Operator == "" {
		fixedNum := stmt.Statement2.Num
		solveStmt = stmt.Statement1

		switch stmt.Operator {
		case "*":
			solveWant = want / fixedNum
		case "+":
			solveWant = want - fixedNum
		case "/":
			solveWant = want * fixedNum
		case "-":
			solveWant = want + fixedNum
		}
	} else {
		panic(fmt.Errorf("No fixed num in statement: %+v", stmt))
	}

	return SolveToResult(solveStmt, solveWant)
}

func SolveForEquality(stmt *Statement) int {
	fixedNum, solveStmt := getFixAndSolve(stmt)
	return SolveToResult(solveStmt, fixedNum)
}

func part2(input string) int {
	statements := parseInput(input)

	statementMap := make(map[string]*Statement, len(statements))
	for _, statement := range statements {
		statementMap[statement.Name] = statement
	}
	statementMap["humn"].Human = true
	statementMap["humn"].Operator = "human"

	stmt := Resolve(statementMap, "root")

	human := SolveForEquality(stmt)

	resetMap(statementMap)
	statementMap["humn"].Human = false
	statementMap["humn"].Operator = ""
	statementMap["humn"].Num = human

	statementMap["root"].Operator = "-"

	test := Resolve(statementMap, "root").Num
	if test != 0 {
		panic(fmt.Errorf("Did not resolve properly: %d", test))
	}

	return human
}

func parseInput(input string) (ans []*Statement) {
	for _, line := range strings.Split(input, "\n") {
		lineSplit := strings.Split(line, ": ")
		statement := &Statement{
			Name: lineSplit[0],
		}

		calcSplit := strings.Split(lineSplit[1], " ")
		if len(calcSplit) == 1 {
			statement.Num = cast.ToInt(calcSplit[0])
		} else {
			statement.Operand1 = calcSplit[0]
			statement.Operator = calcSplit[1]
			statement.Operand2 = calcSplit[2]
		}
		ans = append(ans, statement)
	}
	return ans
}
