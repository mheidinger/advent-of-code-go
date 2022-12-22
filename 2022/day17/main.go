package main

import (
	"crypto/sha256"
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
	DirLeft  = "<"
	DirRight = ">"
)

type Position struct {
	x      int
	height int
}

type Rock interface {
	CheckCollission(chamber [][]bool, newPos Position) bool
	MarkSolid(chamber [][]bool)
	GetPosition() Position
	SetPosition(pos Position)
	GetHeight() int
}

func MoveLeft(chamber [][]bool, rock Rock) {
	pos := rock.GetPosition()
	newPos := Position{height: pos.height, x: pos.x - 1}
	if !rock.CheckCollission(chamber, newPos) {
		rock.SetPosition(newPos)
	}
}

func MoveRight(chamber [][]bool, rock Rock) {
	pos := rock.GetPosition()
	newPos := Position{height: pos.height, x: pos.x + 1}
	if !rock.CheckCollission(chamber, newPos) {
		rock.SetPosition(newPos)
	}
}

func MoveDown(chamber [][]bool, rock Rock) bool {
	pos := rock.GetPosition()
	newPos := Position{height: pos.height - 1, x: pos.x}
	if !rock.CheckCollission(chamber, newPos) {
		rock.SetPosition(newPos)
		return true
	}
	return false
}

func getRowString(row []bool) string {
	builder := strings.Builder{}
	for _, block := range row {
		if block {
			builder.WriteString("#")
		} else {
			builder.WriteString(".")
		}
	}
	return builder.String()
}

func drawChamber(chamber [][]bool, lastHeight int) {
	for height := lastHeight; height >= 0; height-- {
		fmt.Printf("|%s|\n", getRowString(chamber[height]))
	}
	fmt.Println("+-------+")
}

type possibleCycle struct {
	height   int
	numRocks int
}

func runSimulation(input string, numRocks int) int {
	directions := parseInput(input)

	chamberHeight := 10000
	chamberWidth := 7
	chamber := make([][]bool, 0, chamberHeight)
	for it := 0; it < chamberHeight; it++ {
		chamber = append(chamber, make([]bool, chamberWidth))
	}

	possibleCycles := map[string]possibleCycle{}

	rocks := []Rock{&RockHorizontal{}, &RockCross{}, &RockCorner{}, &RockVertical{}, &RockSquare{}}
	rockIt := 0
	directionIt := 0
	maxHeight := -1
	heightOffset := 0
	for it := 0; it < numRocks; it++ {
		currentRock := rocks[rockIt]
		spawnHeight := maxHeight + currentRock.GetHeight() + 3
		currentRock.SetPosition(Position{height: spawnHeight, x: 2})

		for {
			if directions[directionIt] == DirLeft {
				MoveLeft(chamber, currentRock)
			} else {
				MoveRight(chamber, currentRock)
			}
			directionIt = (directionIt + 1) % len(directions)
			if !MoveDown(chamber, currentRock) {
				break
			}
		}

		currentRock.MarkSolid(chamber)
		if currentRock.GetPosition().height > maxHeight {
			maxHeight = currentRock.GetPosition().height
		}

		checkEnd := maxHeight
		checkStart := checkEnd - 10
		if checkStart < 0 {
			checkStart = 0
		}
		hashString := strings.Builder{}
		for it := checkStart; it < checkEnd; it++ {
			hashString.WriteString(getRowString(chamber[it]))
		}
		sha := sha256.Sum256([]byte(hashString.String()))
		key := fmt.Sprintf("%s,%d,%d", sha, rockIt, directionIt)
		if found, ok := possibleCycles[key]; ok {
			cycleNumRocks := it - found.numRocks
			cycleHeight := maxHeight - found.height
			if cycleNumRocks < numRocks-it {
				numCycles := (numRocks - it) / cycleNumRocks
				fmt.Printf("jump ahead from %d with %d cycles of length %d and height %d\n", it, numCycles, cycleNumRocks, cycleHeight)
				it += numCycles * cycleNumRocks
				heightOffset += numCycles * cycleHeight
			}
		} else {
			possibleCycles[key] = possibleCycle{
				numRocks: it,
				height:   heightOffset + maxHeight,
			}
		}

		rockIt = (rockIt + 1) % len(rocks)
	}

	return heightOffset + maxHeight + 1
}

func part1(input string) int {
	return runSimulation(input, 2022)
}

func part2(input string) int {
	return runSimulation(input, 1000000000000)
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		return strings.Split(line, "")
	}
	return ans
}
