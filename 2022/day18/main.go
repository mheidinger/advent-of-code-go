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

type Cube struct {
	x int
	y int
	z int
}

func (cube Cube) Add(cube2 Cube) Cube {
	return Cube{
		cube.x + cube2.x,
		cube.y + cube2.y,
		cube.z + cube2.z,
	}
}

func (cube Cube) Mult(factor int) Cube {
	return Cube{
		cube.x * factor,
		cube.y * factor,
		cube.z * factor,
	}
}

func (cube Cube) toKey() string {
	return fmt.Sprintf("%d,%d,%d", cube.x, cube.y, cube.z)
}

func getSides() []Cube {
	return []Cube{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
	}
}

func getOpenSides(cube Cube, cubes map[string]Cube) int {
	openSides := 0
	for _, side := range getSides() {
		checkCoord := cube.Add(side)
		if _, ok := cubes[checkCoord.toKey()]; !ok {
			openSides++
		}
	}

	return openSides
}

func walkExterior(cube Cube, cubes map[string]Cube, checked map[string]bool, min, max Cube) int {
	// if this coords were already checked, return 0
	if checked[cube.toKey()] {
		return 0
	}
	// if this coords are out of maximum, return 0 to not discover endlessly
	if cube.x < min.x || cube.x > max.x || cube.y < min.y || cube.y > max.y || cube.z < min.z || cube.z > max.z {
		return 0
	}
	// if this is part of the lava, we found one side that is exposed
	// don't check neighbours as we'll only check further on air
	if _, ok := cubes[cube.toKey()]; ok {
		return 1
	}

	// this is air, mark is already as checked to prevent infinite loop
	// neighbour would check us again, we our neighbour, etc.
	checked[cube.toKey()] = true

	// check all neighbour coords for their open sides
	foundCubes := 0
	for _, side := range getSides() {
		foundCubes += walkExterior(cube.Add(side), cubes, checked, min, max)
	}

	return foundCubes
}

func part1(input string) int {
	cubes, _ := parseInput(input)

	openSides := 0
	for _, cube := range cubes {
		openSides += getOpenSides(cube, cubes)
	}

	return openSides
}

func part2(input string) int {
	cubes, max := parseInput(input)

	max = max.Add(Cube{1, 1, 1})
	min := Cube{-1, -1, -1}

	exterior := make(map[string]bool)
	return walkExterior(Cube{0, 0, 0}, cubes, exterior, min, max)
}

func parseInput(input string) (ans map[string]Cube, max Cube) {
	ans = make(map[string]Cube, 0)
	for _, line := range strings.Split(input, "\n") {
		lineSplit := strings.Split(line, ",")
		cube := Cube{
			x: cast.ToInt(lineSplit[0]),
			y: cast.ToInt(lineSplit[1]),
			z: cast.ToInt(lineSplit[2]),
		}
		ans[line] = cube

		if cube.x > max.x {
			max.x = cube.x
		}
		if cube.y > max.y {
			max.y = cube.y
		}
		if cube.z > max.z {
			max.z = cube.z
		}
	}
	return ans, max
}
