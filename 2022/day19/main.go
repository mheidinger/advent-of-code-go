package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/barkimedes/go-deepcopy"
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

type Blueprint struct {
	oreRobotCostOre        int
	clayRobotCostOre       int
	obsidianRobotCostOre   int
	obsidianRobotCostClay  int
	geodeRobotCostOre      int
	geodeRobotCostObsidian int
	maxOreCost             int
	maxClayCost            int
}

type Inventory struct {
	Ores   map[string]int
	Robots map[string]int
}

const (
	MatOre      = "ore"
	MatClay     = "clay"
	MatObsidian = "obsidian"
	MatGeode    = "geode"
)

func (inv Inventory) Copy() Inventory {
	return deepcopy.MustAnything(inv).(Inventory)
}

// return numOre, numClay, numObsidian
func (bp Blueprint) GetCosts(build string) (int, int, int) {
	switch build {
	case MatOre:
		return bp.oreRobotCostOre, 0, 0
	case MatClay:
		return bp.clayRobotCostOre, 0, 0
	case MatObsidian:
		return bp.obsidianRobotCostOre, bp.obsidianRobotCostClay, 0
	case MatGeode:
		return bp.geodeRobotCostOre, 0, bp.geodeRobotCostObsidian
	}
	panic(fmt.Errorf("unknown robot to build: %s", build))
}

func (inv Inventory) TimeUntilReached(requiredOre, requiredClay, requiredObsidian int) int {
	timeUntilOre := 0.0
	timeUntilClay := 0.0
	timeUntilObsidian := 0.0

	if inv.Robots[MatOre] > 0 {
		missingOre := requiredOre - inv.Ores[MatOre]
		timeUntilOre = float64(missingOre) / float64(inv.Robots[MatOre])
	}
	if inv.Robots[MatClay] > 0 {
		missingClay := requiredClay - inv.Ores[MatClay]
		timeUntilClay = float64(missingClay) / float64(inv.Robots[MatClay])
	}
	if inv.Robots[MatObsidian] > 0 {
		missingObsidian := requiredObsidian - inv.Ores[MatObsidian]
		timeUntilObsidian = float64(missingObsidian) / float64(inv.Robots[MatObsidian])
	}

	return int(math.Ceil(MaxFloat(timeUntilOre, timeUntilClay, timeUntilObsidian)))
}

func (inv Inventory) PassTime(time int) {
	for key := range inv.Robots {
		inv.Ores[key] += inv.Robots[key] * time
	}
}

func (bp Blueprint) Simulate(cache map[string]int, inv Inventory, timeLeft int, buildNext string) int {
	keyBuilder := strings.Builder{}
	for _, key := range []string{MatOre, MatClay, MatObsidian, MatGeode} {
		keyBuilder.WriteString(fmt.Sprintf("%d%d", inv.Ores[key], inv.Robots[key]))
	}
	keyBuilder.WriteString(fmt.Sprintf("%d", timeLeft))
	keyBuilder.WriteString(buildNext)
	key := keyBuilder.String()

	if found, ok := cache[key]; ok {
		return found
	}

	if buildNext != "" {
		requiredOre, requiredClay, requiredObsidian := bp.GetCosts(buildNext)
		timeUntilReached := inv.TimeUntilReached(requiredOre, requiredClay, requiredObsidian)
		if timeUntilReached+1 > timeLeft {
			inv.PassTime(timeLeft)
			return inv.Ores[MatGeode]
		}
		inv.PassTime(timeUntilReached + 1)
		inv.Robots[buildNext]++
		inv.Ores[MatOre] -= requiredOre
		inv.Ores[MatClay] -= requiredClay
		inv.Ores[MatObsidian] -= requiredObsidian
		timeLeft -= timeUntilReached + 1
	}

	maxGeodes := 0
	for _, buildNext := range inv.GetPossibleBuilds(bp, timeLeft) {
		geodes := bp.Simulate(cache, inv.Copy(), timeLeft, buildNext)
		if geodes > maxGeodes {
			maxGeodes = geodes
		}
	}

	cache[key] = maxGeodes

	return maxGeodes
}

func (inv Inventory) GetPossibleBuilds(bp Blueprint, timeLeft int) []string {
	builds := []string{}

	if inv.Robots[MatOre]*timeLeft+inv.Ores[MatOre] < bp.maxOreCost*timeLeft {
		builds = append(builds, MatOre)
	}

	if inv.Robots[MatClay]*timeLeft+inv.Ores[MatClay] < bp.maxClayCost*timeLeft {
		builds = append(builds, MatClay)
	}

	if inv.Robots[MatClay] > 0 && inv.Robots[MatObsidian]*timeLeft+inv.Ores[MatObsidian] < bp.geodeRobotCostObsidian*timeLeft {
		builds = append(builds, MatObsidian)
	}

	if inv.Robots[MatObsidian] > 0 {
		builds = append(builds, MatGeode)
	}

	return builds
}

func part1(input string) int {
	blueprints := parseInput(input)

	sumQualityLevel := 0
	for it, bp := range blueprints {
		initialInv := Inventory{
			Ores:   map[string]int{MatOre: 0, MatClay: 0, MatObsidian: 0, MatGeode: 0},
			Robots: map[string]int{MatOre: 1, MatClay: 0, MatObsidian: 0, MatGeode: 0},
		}
		totalTime := 24
		cache := map[string]int{}
		geodes := bp.Simulate(cache, initialInv, totalTime, "")
		fmt.Printf("Blueprint #%d: %d geodes\n", it+1, geodes)
		sumQualityLevel += geodes * (it + 1)
	}
	return sumQualityLevel
}

func part2(input string) int {
	blueprints := parseInput(input)

	prodGeodes := 1
	for it, bp := range blueprints[:3] {
		initialInv := Inventory{
			Ores:   map[string]int{MatOre: 0, MatClay: 0, MatObsidian: 0, MatGeode: 0},
			Robots: map[string]int{MatOre: 1, MatClay: 0, MatObsidian: 0, MatGeode: 0},
		}
		totalTime := 32
		cache := map[string]int{}
		geodes := bp.Simulate(cache, initialInv, totalTime, "")
		fmt.Printf("Blueprint #%d: %d geodes\n", it+1, geodes)
		prodGeodes *= geodes
	}
	return prodGeodes
}

var reg = regexp.MustCompile(`Blueprint \d+: Each ore robot costs (\d+) ore\. Each clay robot costs (\d+) ore\. Each obsidian robot costs (\d+) ore and (\d+) clay\. Each geode robot costs (\d+) ore and (\d+) obsidian\.`)

func Max(nums ...int) int {
	max := 0
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}

func MaxFloat(nums ...float64) float64 {
	max := 0.0
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}

func parseInput(input string) (ans []Blueprint) {
	for _, line := range strings.Split(input, "\n") {
		matches := reg.FindStringSubmatch(line)
		bp := Blueprint{
			cast.ToInt(matches[1]),
			cast.ToInt(matches[2]),
			cast.ToInt(matches[3]),
			cast.ToInt(matches[4]),
			cast.ToInt(matches[5]),
			cast.ToInt(matches[6]),
			0,
			0,
		}
		bp.maxOreCost = Max(bp.clayRobotCostOre, bp.obsidianRobotCostOre, bp.oreRobotCostOre)
		bp.maxClayCost = Max(bp.obsidianRobotCostClay)
		ans = append(ans, bp)
	}
	return ans
}
