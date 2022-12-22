package main

import (
	_ "embed"
	"flag"
	"fmt"
	"sort"
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

type Operand string

const (
	OpMult = "*"
	OpAdd  = "+"
)

type Operation struct {
	op  Operand
	val int
}

func (op Operation) calc(num int) int {
	if op.val == 0 {
		op.val = num
	}
	switch op.op {
	case OpMult:
		return op.val * num
	case OpAdd:
		return op.val + num
	}
	panic(fmt.Errorf("unknown operand: %v", op.op))
}

type Monkey struct {
	items            []int
	op               Operation
	testDivisible    int
	testOkTarget     int
	testFailTarget   int
	totalInspections int
}

type Throw struct {
	item   int
	target int
}

func (m *Monkey) getThrows(divide bool, lcm int) []Throw {
	throws := []Throw{}
	m.totalInspections += len(m.items)
	for _, item := range m.items {
		worryAfterOp := m.op.calc(item)
		worry := worryAfterOp
		if divide {
			worry = int(float64(worryAfterOp) / 3.0)
		} else {
			worry = worry % lcm
		}
		if worry%m.testDivisible == 0 {
			throws = append(throws, Throw{worry, m.testOkTarget})
		} else {
			throws = append(throws, Throw{worry, m.testFailTarget})
		}
	}
	return throws
}

func part1(input string) int {
	monkeys := parseInput(input)

	for it := 0; it < 20; it++ {
		for _, monkey := range monkeys {
			throws := monkey.getThrows(true, 0)
			monkey.items = []int{}
			for _, throw := range throws {
				monkeys[throw.target].items = append(monkeys[throw.target].items, throw.item)
			}
		}
	}

	inspections := []int{}
	for it, monkey := range monkeys {
		fmt.Printf("#%d: %d\n", it, monkey.totalInspections)
		inspections = append(inspections, monkey.totalInspections)
	}
	sort.Ints(inspections)

	return inspections[len(inspections)-1] * inspections[len(inspections)-2]
}

func part2(input string) int {
	monkeys := parseInput(input)

	lcm := 1
	for _, monkey := range monkeys {
		lcm *= monkey.testDivisible
	}

	for it := 0; it < 10000; it++ {
		for _, monkey := range monkeys {
			throws := monkey.getThrows(false, lcm)
			monkey.items = []int{}
			for _, throw := range throws {
				monkeys[throw.target].items = append(monkeys[throw.target].items, throw.item)
			}
		}
	}

	inspections := []int{}
	for it, monkey := range monkeys {
		fmt.Printf("#%d: %d\n", it, monkey.totalInspections)
		inspections = append(inspections, monkey.totalInspections)
	}
	sort.Ints(inspections)

	return inspections[len(inspections)-1] * inspections[len(inspections)-2]
}

func parseInput(input string) (ans []*Monkey) {
	lines := strings.Split(input, "\n")

	for it := 0; it < 8; it++ {
		monk := &Monkey{items: []int{}}
		itemsList := strings.Split(lines[it*7+1], ": ")[1]
		items := strings.Split(itemsList, ", ")
		for _, item := range items {
			monk.items = append(monk.items, cast.ToInt(item))
		}
		opLineSplit := strings.Split(strings.TrimSpace(lines[it*7+2]), " ")
		monk.op = Operation{op: Operand(opLineSplit[4])}
		if opLineSplit[5] == "old" {
			monk.op.val = 0
		} else {
			monk.op.val = cast.ToInt(opLineSplit[5])
		}
		testStr := strings.Split(lines[it*7+3], "by ")[1]
		monk.testDivisible = cast.ToInt(testStr)
		trueStr := strings.Split(lines[it*7+4], "monkey ")[1]
		monk.testOkTarget = cast.ToInt(trueStr)
		falseStr := strings.Split(lines[it*7+5], "monkey ")[1]
		monk.testFailTarget = cast.ToInt(falseStr)
		ans = append(ans, monk)
	}
	return ans
}
