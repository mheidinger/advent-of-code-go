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

type Item struct {
	parent  *Item
	content map[string]*Item
	size    int
	name    string
}

func (i *Item) Print(indent int) {
	text := ""
	for it := 0; it < indent; it++ {
		text += "  "
	}
	text += "- " + i.name + " ("
	if len(i.content) > 0 {
		text += "dir, len="
	} else {
		text += "file, len="
	}
	text += fmt.Sprint(i.size) + ")"

	fmt.Println(text)

	for _, child := range i.content {
		child.Print(indent + 1)
	}
}

func (i *Item) CalcSize() int {
	if i.size > 0 {
		return i.size
	}

	for _, child := range i.content {
		i.size += child.CalcSize()
	}
	return i.size
}

func (i *Item) SumLargeDirsBelowThreshold(threshold int) int {
	if len(i.content) == 0 {
		return 0
	}

	returnSize := 0
	if i.size <= threshold {
		returnSize += i.size
	}
	for _, child := range i.content {
		returnSize += child.SumLargeDirsBelowThreshold(threshold)
	}
	return returnSize
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (i *Item) FindDirToDelete(spaceToFree, currDeletionSize int) int {
	if i.size > spaceToFree && i.size-spaceToFree < currDeletionSize-spaceToFree {
		currDeletionSize = i.size
	}
	for _, child := range i.content {
		childDeletionSize := child.FindDirToDelete(spaceToFree, currDeletionSize)
		if childDeletionSize-spaceToFree < currDeletionSize-spaceToFree {
			currDeletionSize = childDeletionSize
		}
	}
	return currDeletionSize
}

func part1(input string) int {
	parsed := parseInput(input)

	parsed.CalcSize()

	return parsed.SumLargeDirsBelowThreshold(100000)
}

func part2(input string) int {
	parsed := parseInput(input)
	parsed.CalcSize()

	totalSize := 70000000
	spaceNeeded := 30000000

	unusedSpace := totalSize - parsed.size
	spaceToFree := spaceNeeded - unusedSpace

	return parsed.FindDirToDelete(spaceToFree, totalSize)
}

func parseInput(input string) (root *Item) {
	root = &Item{
		name:    "/",
		content: map[string]*Item{},
	}
	currentFolder := root
	for _, line := range strings.Split(input, "\n") {
		if strings.HasPrefix(line, "$ ") {
			cmdParts := strings.Split(line[2:], " ")
			switch cmdParts[0] {
			case "cd":
				if cmdParts[1] == ".." {
					currentFolder = currentFolder.parent
				} else {
					currentFolder = currentFolder.content[cmdParts[1]]
				}
			case "ls":
				// Do nothing, next lines without a dollar will be the content
			default:
				panic(fmt.Errorf("unknown command: %v", cmdParts))
			}
			continue
		}

		// Process content of current folder
		contentSplit := strings.Split(line, " ")
		newItem := &Item{
			name:    contentSplit[1],
			parent:  currentFolder,
			content: map[string]*Item{},
		}
		if contentSplit[0] != "dir" {
			newItem.size = cast.ToInt(contentSplit[0])
		}
		currentFolder.content[newItem.name] = newItem
	}
	return root
}
