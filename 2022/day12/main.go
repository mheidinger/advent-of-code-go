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

const MAX_DISTANCE = 999999

type Point struct {
	x int
	y int
}

type Tile struct {
	pos               Point
	height            int
	visited           bool
	tentativeDistance int
}

func findShortestPath(land [][]*Tile, start, goal Point) {
	current := land[start.x][start.y]

	for current != nil {
		neighbours := getNeighbours(land, current)
		for _, neighbour := range neighbours {
			if current.tentativeDistance+1 < neighbour.tentativeDistance {
				neighbour.tentativeDistance = current.tentativeDistance + 1
			}
		}

		current.visited = true
		if current.pos == goal {
			break
		}
		current = getLowestTile(land)
	}
}

func getNeighbours(land [][]*Tile, current *Tile) (neighbours []*Tile) {
	sides := []Point{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	}

	currentHeight := current.height
	for _, side := range sides {
		checkPoint := Point{current.pos.x + side.x, current.pos.y + side.y}
		if checkPoint.x < 0 || checkPoint.x >= len(land) || checkPoint.y < 0 || checkPoint.y >= len(land[0]) {
			continue
		}
		neighbour := land[checkPoint.x][checkPoint.y]
		if neighbour.visited || neighbour.height > currentHeight+1 {
			continue
		}
		neighbours = append(neighbours, neighbour)
	}
	return
}

func getLowestTile(land [][]*Tile) *Tile {
	var minimum *Tile
	for _, row := range land {
		for _, tile := range row {
			if (minimum == nil || tile.tentativeDistance < minimum.tentativeDistance) && !tile.visited {
				minimum = tile
			}
		}
	}
	return minimum
}

func resetLand(land [][]*Tile, start Point) {
	for _, row := range land {
		for _, tile := range row {
			tile.tentativeDistance = MAX_DISTANCE
			tile.visited = false
		}
	}
	land[start.x][start.y].tentativeDistance = 0
}

func part1(input string) int {
	land, start, goal := parseInput(input)

	resetLand(land, start)
	findShortestPath(land, start, goal)

	return land[goal.x][goal.y].tentativeDistance
}

func part2(input string) int {
	land, _, goal := parseInput(input)

	// Better to switch start and goal, only one calculation!
	shortestPath := MAX_DISTANCE
	for _, row := range land {
		for _, tile := range row {
			if tile.height == 0 {
				resetLand(land, tile.pos)
				findShortestPath(land, tile.pos, goal)
				if land[goal.x][goal.y].tentativeDistance < shortestPath {
					shortestPath = land[goal.x][goal.y].tentativeDistance
				}
			}
		}
	}

	return shortestPath
}

func parseInput(input string) (ans [][]*Tile, start, goal Point) {
	for x, line := range strings.Split(input, "\n") {
		lineSplit := strings.Split(line, "")
		row := []*Tile{}
		for y, char := range lineSplit {
			tile := &Tile{
				pos:     Point{x, y},
				visited: false,
			}
			if char == "S" {
				start.x = x
				start.y = y
				char = "a"
			} else if char == "E" {
				goal.x = x
				goal.y = y
				char = "z"
			}
			tile.height = cast.ToASCIICode(char) - cast.ASCIICodeLowerA
			row = append(row, tile)
		}
		ans = append(ans, row)
	}
	return ans, start, goal
}
