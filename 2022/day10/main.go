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

type Operator string

const (
	OpNoop = "noop"
	OpAddx = "addx"
)

type Instruction struct {
	op  Operator
	val int
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

func part1(input string) int {
	instructions := parseInput(input)
	sigStrengthSum := 0

	instrIt := 0
	x := 1
	secondCycle := false
	for cycle := 1; ; cycle++ {
		if cycle == 20 || cycle == 60 || cycle == 100 || cycle == 140 || cycle == 180 || cycle == 220 {
			sigStrengthSum += cycle * x
		}
		instr := instructions[instrIt]
		switch instr.op {
		case OpNoop:
			instrIt++
		case OpAddx:
			if secondCycle {
				instrIt++
				x += instr.val
				secondCycle = false
			} else {
				secondCycle = true
			}
		}

		if instrIt >= len(instructions) {
			break
		}
	}

	return sigStrengthSum
}

func part2(input string) int {
	instructions := parseInput(input)

	instrIt := 0
	x := 1
	secondCycle := false
	for cycle := 0; ; cycle++ {
		crtX := cycle % 40
		if crtX == 0 {
			fmt.Println()
		}
		if x >= crtX-1 && x <= crtX+1 {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}

		instr := instructions[instrIt]
		switch instr.op {
		case OpNoop:
			instrIt++
		case OpAddx:
			if secondCycle {
				instrIt++
				x += instr.val
				secondCycle = false
			} else {
				secondCycle = true
			}
		}

		if instrIt >= len(instructions) {
			break
		}
	}
	fmt.Println()

	return 0
}

func parseInput(input string) (ans []Instruction) {
	for _, line := range strings.Split(input, "\n") {
		instrSplit := strings.Split(line, " ")
		instr := Instruction{OpNoop, 0}
		if instrSplit[0] == OpAddx {
			instr.op = OpAddx
			instr.val = cast.ToInt(instrSplit[1])
		}
		ans = append(ans, instr)
	}
	return ans
}
