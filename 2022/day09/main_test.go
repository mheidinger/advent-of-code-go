package main

import (
	"testing"
)

func TestFollow(t *testing.T) {
	tests := []struct {
		name         string
		inputHead    Position
		inputTail    Position
		expectedTail Position
	}{
		{
			name:         "same",
			inputHead:    Position{0, 0},
			inputTail:    Position{0, 0},
			expectedTail: Position{0, 0},
		},
		{
			name:         "1up",
			inputHead:    Position{1, 0},
			inputTail:    Position{0, 0},
			expectedTail: Position{0, 0},
		},
		{
			name:         "1right",
			inputHead:    Position{0, 1},
			inputTail:    Position{0, 0},
			expectedTail: Position{0, 0},
		},
		{
			name:         "2up",
			inputHead:    Position{2, 0},
			inputTail:    Position{0, 0},
			expectedTail: Position{1, 0},
		},
		{
			name:         "2right",
			inputHead:    Position{0, 2},
			inputTail:    Position{0, 0},
			expectedTail: Position{0, 1},
		},
		{
			name:         "2up1right",
			inputHead:    Position{2, 1},
			inputTail:    Position{0, 0},
			expectedTail: Position{1, 1},
		},
		{
			name:         "2down1left",
			inputHead:    Position{-2, -1},
			inputTail:    Position{0, 0},
			expectedTail: Position{-1, -1},
		},
		{
			name:         "2right1up",
			inputHead:    Position{1, 2},
			inputTail:    Position{0, 0},
			expectedTail: Position{1, 1},
		},
		{
			name:         "2right2up",
			inputHead:    Position{2, 2},
			inputTail:    Position{0, 0},
			expectedTail: Position{1, 1},
		},
		{
			name:         "2left2down",
			inputHead:    Position{-2, -2},
			inputTail:    Position{0, 0},
			expectedTail: Position{-1, -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := follow(tt.inputHead, tt.inputTail); got != tt.expectedTail {
				t.Errorf("follow() = %v, want %v", got, tt.expectedTail)
			}
		})
	}
}

var example = `R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2`

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example",
			input: example,
			want:  13,
		},
		// {
		// 	name:  "actual",
		// 	input: input,
		// 	want:  0,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.input); got != tt.want {
				t.Errorf("part1() = %v, want %v", got, tt.want)
			}
		})
	}
}

var example2 = `R 5
U 8
L 8
D 3
R 17
D 10
L 25
U 20`

func Test_part2(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example",
			input: example,
			want:  1,
		},
		{
			name:  "example2",
			input: example2,
			want:  36,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part2(tt.input); got != tt.want {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}
		})
	}
}
