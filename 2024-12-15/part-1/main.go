package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func ParseInputFile(path string) ([][]byte, []byte) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		panic(0)
	}
	defer file.Close()

	var warehouse [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) <= 0 {
			break
		}
		warehouse = append(warehouse, []byte(line))
	}

	var instructions []byte
	for scanner.Scan() {
		line := scanner.Text()
		instructions = append(instructions, []byte(line)...)
	}

	return warehouse, instructions
}

func main() {
	var path string
	flag.StringVar(&path, "path", "", "The path to the input file")

	var debug bool
	flag.BoolVar(&debug, "debug", false, "Specify whether or not to produce debug output")

	flag.Parse()

	warehouse, instructions := ParseInputFile(path)

	if debug {
		fmt.Println("Starting state of warehouse")
		PrintWarehouse(warehouse)
	}

	var row int
	var col int
	var foundRobot bool = false
	for i, items := range warehouse {
		for j, item := range items {
			if item == '@' {
				row = i
				col = j
				foundRobot = true
				break
			}
		}

		if foundRobot {
			break
		}
	}

	robot := &Robot{
		row: row,
		col: col,
	}

	for _, instruction := range instructions {
		robot.Execute(warehouse, instruction)

		if debug {
			fmt.Printf("\nInstruction: %q\n", instruction)
			PrintWarehouse(warehouse)
		}
	}

	fmt.Println("Sum of all boxes' GPS coordinates: ", SumBoxGpsCoordinates(warehouse))
}

func PrintWarehouse(warehouse [][]byte) {
	fmt.Print("\n")
	for _, items := range warehouse {
		fmt.Println(string(items[:]))
	}
	fmt.Print("\n")
}

type Robot struct {
	row int
	col int
}

func (r *Robot) Execute(warehouse [][]byte, instruction byte) {
	if instruction == '<' {
		if r.col == 0 {
			return
		}

		for col := r.col - 1; col >= 0; col-- {
			item := warehouse[r.row][col]
			if item == '#' {
				return
			} else if item == '.' {
				warehouse[r.row][col] = 'O'
				warehouse[r.row][r.col-1] = '@'
				warehouse[r.row][r.col] = '.'
				r.col -= 1
				return
			}
		}

	} else if instruction == '^' {
		if r.row == 0 {
			return
		}

		for row := r.row - 1; row >= 0; row-- {
			item := warehouse[row][r.col]
			if item == '#' {
				return
			} else if item == '.' {
				warehouse[row][r.col] = 'O'
				warehouse[r.row-1][r.col] = '@'
				warehouse[r.row][r.col] = '.'
				r.row -= 1
				return
			}
		}

	} else if instruction == '>' {
		maxCol := len(warehouse[0]) - 1
		if r.col == maxCol {
			return
		}

		for col := r.col + 1; col <= maxCol; col++ {
			item := warehouse[r.row][col]
			if item == '#' {
				return
			} else if item == '.' {
				warehouse[r.row][col] = 'O'
				warehouse[r.row][r.col+1] = '@'
				warehouse[r.row][r.col] = '.'
				r.col += 1
				return
			}
		}

	} else if instruction == 'v' {
		maxRow := len(warehouse) - 1
		if r.row == maxRow {
			return
		}

		for row := r.row + 1; row <= maxRow; row++ {
			item := warehouse[row][r.col]
			if item == '#' {
				return
			} else if item == '.' {
				warehouse[row][r.col] = 'O'
				warehouse[r.row+1][r.col] = '@'
				warehouse[r.row][r.col] = '.'
				r.row += 1
				return
			}
		}
	}
}

func SumBoxGpsCoordinates(warehouse [][]byte) int {
	runningTotal := 0
	for i, row := range warehouse {
		for j, item := range row {
			if item == 'O' {
				runningTotal += 100*i + j
			}
		}
	}

	return runningTotal
}
