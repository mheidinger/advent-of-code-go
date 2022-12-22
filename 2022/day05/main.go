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

type Command struct {
	Amount int
	From   int
	To     int
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

func reverse(s []string) []string {
	a := make([]string, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}

func part1(input string) string {
	commands := parseInput(input)
	stacks := getInitialState()

	for _, cmd := range commands {
		fromLen := len(stacks[cmd.From])
		move := reverse(stacks[cmd.From][fromLen-cmd.Amount:])
		stacks[cmd.From] = stacks[cmd.From][:fromLen-cmd.Amount]
		stacks[cmd.To] = append(stacks[cmd.To], move...)
	}

	res := ""
	for _, stack := range stacks {
		res += stack[len(stack)-1]
	}

	return res
}

func part2(input string) string {
	commands := parseInput(input)
	stacks := getInitialState()

	for _, cmd := range commands {
		fromLen := len(stacks[cmd.From])
		move := stacks[cmd.From][fromLen-cmd.Amount:]
		stacks[cmd.From] = stacks[cmd.From][:fromLen-cmd.Amount]
		stacks[cmd.To] = append(stacks[cmd.To], move...)
	}

	res := ""
	for _, stack := range stacks {
		res += stack[len(stack)-1]
	}

	return res
}

/*
[N]     [Q]         [N]
[R]     [F] [Q]     [G] [M]
[J]     [Z] [T]     [R] [H] [J]
[T] [H] [G] [R]     [B] [N] [T]
[Z] [J] [J] [G] [F] [Z] [S] [M]
[B] [N] [N] [N] [Q] [W] [L] [Q] [S]
[D] [S] [R] [V] [T] [C] [C] [N] [G]
[F] [R] [C] [F] [L] [Q] [F] [D] [P]
 1   2   3   4   5   6   7   8   9

 => Shifted to 0-8
*/

func getInitialState() [][]string {
	stacks := [][]string{}
	stacks = append(stacks, []string{"F", "D", "B", "Z", "T", "J", "R", "N"})
	stacks = append(stacks, []string{"R", "S", "N", "J", "H"})
	stacks = append(stacks, []string{"C", "R", "N", "J", "G", "Z", "F", "Q"})
	stacks = append(stacks, []string{"F", "V", "N", "G", "R", "T", "Q"})
	stacks = append(stacks, []string{"L", "T", "Q", "F"})
	stacks = append(stacks, []string{"Q", "C", "W", "Z", "B", "R", "G", "N"})
	stacks = append(stacks, []string{"F", "C", "L", "S", "N", "H", "M"})
	stacks = append(stacks, []string{"D", "N", "Q", "M", "T", "J"})
	stacks = append(stacks, []string{"P", "G", "S"})
	return stacks
}

func getInitialTestState() [][]string {
	stacks := [][]string{}
	stacks = append(stacks, []string{"Z", "N"})
	stacks = append(stacks, []string{"M", "C", "D"})
	stacks = append(stacks, []string{"P"})
	return stacks
}

func parseInput(input string) (ans []Command) {
	for _, line := range strings.Split(input, "\n") {
		lineSplit := strings.Split(line, " ")
		cmd := Command{
			Amount: cast.ToInt(lineSplit[1]),
			From:   cast.ToInt(lineSplit[3]) - 1,
			To:     cast.ToInt(lineSplit[5]) - 1,
		}
		ans = append(ans, cmd)
	}
	return ans
}
