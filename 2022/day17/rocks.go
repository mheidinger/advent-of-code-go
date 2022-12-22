package main

// ####

type RockHorizontal struct {
	Position
}

var _ Rock = &RockHorizontal{}

func (rock *RockHorizontal) GetHeight() int {
	return 1
}

func (rock *RockHorizontal) SetPosition(pos Position) {
	rock.Position = pos
}

func (rock *RockHorizontal) GetPosition() Position {
	return rock.Position
}

func (rock *RockHorizontal) CheckCollission(chamber [][]bool, newPos Position) (collision bool) {
	defer func() {
		if err := recover(); err != nil {
			collision = true
		}
	}()

	if !chamber[newPos.height][newPos.x] &&
		!chamber[newPos.height][newPos.x+1] &&
		!chamber[newPos.height][newPos.x+2] &&
		!chamber[newPos.height][newPos.x+3] {
		collision = false
	} else {
		collision = true
	}
	return
}

func (rock *RockHorizontal) MarkSolid(chamber [][]bool) {
	chamber[rock.height][rock.x] = true
	chamber[rock.height][rock.x+1] = true
	chamber[rock.height][rock.x+2] = true
	chamber[rock.height][rock.x+3] = true
}

// .#.
// ###
// .#.

type RockCross struct {
	Position
}

var _ Rock = &RockCross{}

func (rock *RockCross) GetHeight() int {
	return 3
}

func (rock *RockCross) SetPosition(pos Position) {
	rock.Position = pos
}

func (rock *RockCross) GetPosition() Position {
	return rock.Position
}

func (rock *RockCross) CheckCollission(chamber [][]bool, newPos Position) (collision bool) {
	defer func() {
		if err := recover(); err != nil {
			collision = true
		}
	}()

	if !chamber[newPos.height][newPos.x+1] &&
		!chamber[newPos.height-1][newPos.x] &&
		!chamber[newPos.height-1][newPos.x+1] &&
		!chamber[newPos.height-1][newPos.x+2] &&
		!chamber[newPos.height-2][newPos.x+1] {
		collision = false
	} else {
		collision = true
	}
	return
}

func (rock *RockCross) MarkSolid(chamber [][]bool) {
	chamber[rock.height][rock.x+1] = true
	chamber[rock.height-1][rock.x] = true
	chamber[rock.height-1][rock.x+1] = true
	chamber[rock.height-1][rock.x+2] = true
	chamber[rock.height-2][rock.x+1] = true
}

// ..#
// ..#
// ###

type RockCorner struct {
	Position
}

var _ Rock = &RockCorner{}

func (rock *RockCorner) GetHeight() int {
	return 3
}

func (rock *RockCorner) SetPosition(pos Position) {
	rock.Position = pos
}

func (rock *RockCorner) GetPosition() Position {
	return rock.Position
}

func (rock *RockCorner) CheckCollission(chamber [][]bool, newPos Position) (collision bool) {
	defer func() {
		if err := recover(); err != nil {
			collision = true
		}
	}()

	if !chamber[newPos.height][newPos.x+2] &&
		!chamber[newPos.height-1][newPos.x+2] &&
		!chamber[newPos.height-2][newPos.x] &&
		!chamber[newPos.height-2][newPos.x+1] &&
		!chamber[newPos.height-2][newPos.x+2] {
		collision = false
	} else {
		collision = true
	}
	return
}

func (rock *RockCorner) MarkSolid(chamber [][]bool) {
	chamber[rock.height][rock.x+2] = true
	chamber[rock.height-1][rock.x+2] = true
	chamber[rock.height-2][rock.x] = true
	chamber[rock.height-2][rock.x+1] = true
	chamber[rock.height-2][rock.x+2] = true
}

// #
// #
// #
// #

type RockVertical struct {
	Position
}

var _ Rock = &RockVertical{}

func (rock *RockVertical) GetHeight() int {
	return 4
}

func (rock *RockVertical) SetPosition(pos Position) {
	rock.Position = pos
}

func (rock *RockVertical) GetPosition() Position {
	return rock.Position
}

func (rock *RockVertical) CheckCollission(chamber [][]bool, newPos Position) (collision bool) {
	defer func() {
		if err := recover(); err != nil {
			collision = true
		}
	}()

	if !chamber[newPos.height][newPos.x] &&
		!chamber[newPos.height-1][newPos.x] &&
		!chamber[newPos.height-2][newPos.x] &&
		!chamber[newPos.height-3][newPos.x] {
		collision = false
	} else {
		collision = true
	}
	return
}

func (rock *RockVertical) MarkSolid(chamber [][]bool) {
	chamber[rock.height][rock.x] = true
	chamber[rock.height-1][rock.x] = true
	chamber[rock.height-2][rock.x] = true
	chamber[rock.height-3][rock.x] = true
}

// ##
// ##

type RockSquare struct {
	Position
}

var _ Rock = &RockSquare{}

func (rock *RockSquare) GetHeight() int {
	return 2
}

func (rock *RockSquare) SetPosition(pos Position) {
	rock.Position = pos
}

func (rock *RockSquare) GetPosition() Position {
	return rock.Position
}

func (rock *RockSquare) CheckCollission(chamber [][]bool, newPos Position) (collision bool) {
	defer func() {
		if err := recover(); err != nil {
			collision = true
		}
	}()

	if !chamber[newPos.height][newPos.x] &&
		!chamber[newPos.height][newPos.x+1] &&
		!chamber[newPos.height-1][newPos.x] &&
		!chamber[newPos.height-1][newPos.x+1] {
		collision = false
	} else {
		collision = true
	}
	return
}

func (rock *RockSquare) MarkSolid(chamber [][]bool) {
	chamber[rock.height][rock.x] = true
	chamber[rock.height][rock.x+1] = true
	chamber[rock.height-1][rock.x] = true
	chamber[rock.height-1][rock.x+1] = true
}
