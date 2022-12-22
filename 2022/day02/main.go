package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

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

const (
	scoreLoss = 0
	scoreDraw = 3
	scoreWin  = 6
)

type Play string

const (
	PlayRock    = "rock"
	PlayPaper   = "paper"
	PlayScissor = "scissor"
)

var playScores = map[Play]int{
	PlayRock:    1,
	PlayPaper:   2,
	PlayScissor: 3,
}

type Outcome string

const (
	OutcomeLoss = "X"
	OutcomeDraw = "Y"
	OutcomeWin  = "Z"
)

var outcomeScores = map[Outcome]int{
	OutcomeLoss: 0,
	OutcomeDraw: 3,
	OutcomeWin:  6,
}

type Round struct {
	Opponent        Play
	Self            Play
	ExpectedOutcome Outcome
}

func (round Round) GetSelfScore() int {
	var outcome Outcome
	if round.Self == round.Opponent {
		outcome = OutcomeDraw
	} else if round.Self == PlayRock && round.Opponent == PlayScissor {
		outcome = OutcomeWin
	} else if round.Self == PlayPaper && round.Opponent == PlayScissor {
		outcome = OutcomeLoss
	} else if round.Self == PlayRock && round.Opponent == PlayPaper {
		outcome = OutcomeLoss
	} else if round.Self == PlayPaper && round.Opponent == PlayRock {
		outcome = OutcomeWin
	} else if round.Self == PlayScissor && round.Opponent == PlayPaper {
		outcome = OutcomeWin
	} else if round.Self == PlayScissor && round.Opponent == PlayRock {
		outcome = OutcomeLoss
	} else {
		panic(fmt.Errorf("not handled situation: '%s' vs '%s'", round.Self, round.Opponent))
	}
	return playScores[round.Self] + outcomeScores[outcome]
}

func (round Round) GetExpectedOutcomeScore() int {
	var play Play
	if round.ExpectedOutcome == OutcomeDraw {
		play = round.Opponent
	} else if round.ExpectedOutcome == OutcomeWin {
		switch round.Opponent {
		case PlayPaper:
			play = PlayScissor
		case PlayRock:
			play = PlayPaper
		case PlayScissor:
			play = PlayRock
		}
	} else if round.ExpectedOutcome == OutcomeLoss {
		switch round.Opponent {
		case PlayPaper:
			play = PlayRock
		case PlayRock:
			play = PlayScissor
		case PlayScissor:
			play = PlayPaper
		}
	}
	return outcomeScores[round.ExpectedOutcome] + playScores[play]
}

func part1(input string) int {
	parsed := parseInputPart1(input)

	totalScore := 0
	for _, round := range parsed {
		totalScore += round.GetSelfScore()
	}

	return totalScore
}

func part2(input string) int {
	parsed := parseInputPart2(input)

	totalScore := 0
	for _, round := range parsed {
		totalScore += round.GetExpectedOutcomeScore()
	}

	return totalScore
}

func parseInputPart1(input string) (ans []Round) {
	inputOpponent := map[string]string{
		"A": PlayRock,
		"B": PlayPaper,
		"C": PlayScissor,
	}

	inputSelf := map[string]string{
		"X": PlayRock,
		"Y": PlayPaper,
		"Z": PlayScissor,
	}

	for _, line := range strings.Split(input, "\n") {
		splits := strings.Split(line, " ")
		ans = append(ans, Round{Opponent: Play(inputOpponent[splits[0]]), Self: Play(inputSelf[splits[1]])})
	}
	return ans
}

func parseInputPart2(input string) (ans []Round) {
	inputOpponent := map[string]string{
		"A": PlayRock,
		"B": PlayPaper,
		"C": PlayScissor,
	}

	for _, line := range strings.Split(input, "\n") {
		splits := strings.Split(line, " ")
		ans = append(ans, Round{Opponent: Play(inputOpponent[splits[0]]), ExpectedOutcome: Outcome(splits[1])})
	}
	return ans
}
