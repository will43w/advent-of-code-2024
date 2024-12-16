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

	oldWarehouse, instructions := ParseInputFile(path)
	newWarehouse := ResizeWarehouse(oldWarehouse)

	if debug {
		fmt.Println("Initial state: ")
		PrintWarehouse(newWarehouse)
	}

	var row int
	var col int
	var foundRobot bool = false
	for i, items := range newWarehouse {
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
		robot.Execute(newWarehouse, instruction)

		if debug {
			fmt.Printf("\nMove: %q:\n", instruction)
			PrintWarehouse(newWarehouse)
		}
	}

	fmt.Println("Sum of all boxes' GPS coordinates: ", SumBoxGpsCoordinates(newWarehouse))
}

func ResizeWarehouse(warehouse [][]byte) [][]byte {
	newWarehouse := [][]byte{{}}

	numCols := len(warehouse[0]) * 2
	for _, items := range warehouse {
		newWarehouseItems := make([]byte, numCols)

		for col, item := range items {
			if item == '#' {
				newWarehouseItems[col*2] = '#'
				newWarehouseItems[col*2+1] = '#'
			} else if item == 'O' {
				newWarehouseItems[col*2] = '['
				newWarehouseItems[col*2+1] = ']'
			} else if item == '.' {
				newWarehouseItems[col*2] = '.'
				newWarehouseItems[col*2+1] = '.'
			} else if item == '@' {
				newWarehouseItems[col*2] = '@'
				newWarehouseItems[col*2+1] = '.'
			}
		}

		newWarehouse = append(newWarehouse, newWarehouseItems)
	}

	return newWarehouse
}

func PrintWarehouse(warehouse [][]byte) {
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
		if warehouse[r.row][r.col-1] == '#' {
			// # @
			return
		}

		if warehouse[r.row][r.col-1] == '.' {
			// . @
			warehouse[r.row][r.col-1] = '@'
			warehouse[r.row][r.col] = '.'
			r.col -= 1
			return
		}

		// [ ] @
		if CanPushBoxLeft(warehouse, r.row, r.col-1) {
			PushBoxLeft(warehouse, r.row, r.col-1)
			warehouse[r.row][r.col-1] = '@'
			warehouse[r.row][r.col] = '.'
			r.col -= 1
		}

	} else if instruction == '^' {
		if warehouse[r.row-1][r.col] == '#' {
			// #
			// @
			return
		}

		if warehouse[r.row-1][r.col] == '.' {
			// .
			// @
			warehouse[r.row-1][r.col] = '@'
			warehouse[r.row][r.col] = '.'
			r.row -= 1
			return
		}

		// [ ] or [ ]
		// @        @
		var leftSideCol int
		itemAboveRobot := warehouse[r.row-1][r.col]
		if itemAboveRobot == '[' {
			leftSideCol = r.col
		} else {
			leftSideCol = r.col - 1
		}
		if CanPushBoxUp(warehouse, r.row-1, leftSideCol) {
			PushBoxUp(warehouse, r.row-1, leftSideCol)
			warehouse[r.row-1][r.col] = '@'
			warehouse[r.row][r.col] = '.'
			r.row -= 1
		}

	} else if instruction == '>' {
		if warehouse[r.row][r.col+1] == '#' {
			// @ #
			return
		}

		if warehouse[r.row][r.col+1] == '.' {
			// @ .
			warehouse[r.row][r.col+1] = '@'
			warehouse[r.row][r.col] = '.'
			r.col += 1
			return
		}

		// @ [ ]
		if CanPushBoxRight(warehouse, r.row, r.col+1) {
			PushBoxRight(warehouse, r.row, r.col+1)
			warehouse[r.row][r.col+1] = '@'
			warehouse[r.row][r.col] = '.'
			r.col += 1
		}

	} else if instruction == 'v' {
		if warehouse[r.row+1][r.col] == '#' {
			// @
			// #
			return
		}

		if warehouse[r.row+1][r.col] == '.' {
			// @
			// .
			warehouse[r.row+1][r.col] = '@'
			warehouse[r.row][r.col] = '.'
			r.row += 1
			return
		}

		// @   or   @
		// [ ]    [ ]
		var leftSideCol int
		itemBelowRobot := warehouse[r.row+1][r.col]
		if itemBelowRobot == '[' {
			leftSideCol = r.col
		} else {
			leftSideCol = r.col - 1
		}
		if CanPushBoxDown(warehouse, r.row+1, leftSideCol) {
			PushBoxDown(warehouse, r.row+1, leftSideCol)
			warehouse[r.row+1][r.col] = '@'
			warehouse[r.row][r.col] = '.'
			r.row += 1
		}
	}
}

func CanPushBoxLeft(warehouse [][]byte, row int, rightSideCol int) bool {
	// # [ ]
	if warehouse[row][rightSideCol-2] == '#' {
		return false
	}

	// . [ ]
	if warehouse[row][rightSideCol-2] == '.' {
		return true
	}

	// [ ] [ ]
	if warehouse[row][rightSideCol-2] == ']' {
		return CanPushBoxLeft(warehouse, row, rightSideCol-2)
	}

	return false
}

func PushBoxLeft(warehouse [][]byte, row int, rightSideCol int) {
	if warehouse[row][rightSideCol-2] == ']' {
		PushBoxLeft(warehouse, row, rightSideCol-2)
	}

	warehouse[row][rightSideCol-2] = '['
	warehouse[row][rightSideCol-1] = ']'
	warehouse[row][rightSideCol] = '.'
}

func CanPushBoxRight(warehouse [][]byte, row int, leftSideCol int) bool {
	// [ ] #
	if warehouse[row][leftSideCol+2] == '#' {
		return false
	}

	// [ ] .
	if warehouse[row][leftSideCol+2] == '.' {
		return true
	}

	// [ ] [ ]
	if warehouse[row][leftSideCol+2] == '[' {
		return CanPushBoxRight(warehouse, row, leftSideCol+2)
	}

	return false
}

func PushBoxRight(warehouse [][]byte, row int, leftSideCol int) {
	if warehouse[row][leftSideCol+2] == '[' {
		PushBoxRight(warehouse, row, leftSideCol+2)
	}

	warehouse[row][leftSideCol+1] = '['
	warehouse[row][leftSideCol+2] = ']'
	warehouse[row][leftSideCol] = '.'
}

func CanPushBoxUp(warehouse [][]byte, row int, leftSideCol int) bool {
	itemAboveLeftSide := warehouse[row-1][leftSideCol]
	itemAboveRightSide := warehouse[row-1][leftSideCol+1]

	// # #  or  # .  or  . #
	// [ ]      [ ]      [ ]
	if itemAboveLeftSide == '#' || itemAboveRightSide == '#' {
		return false
	}

	// . .
	// [ ]
	if itemAboveLeftSide == '.' && itemAboveRightSide == '.' {
		return true
	}

	if itemAboveLeftSide == '[' {
		// [ ]
		// [ ]
		return CanPushBoxUp(warehouse, row-1, leftSideCol)
	} else if itemAboveLeftSide == ']' {
		if itemAboveRightSide == '[' {
			// [ ] [ ]
			//   [ ]
			return CanPushBoxUp(warehouse, row-1, leftSideCol-1) && CanPushBoxUp(warehouse, row-1, leftSideCol+1)
		} else {
			// [ ] .
			//   [ ]
			return CanPushBoxUp(warehouse, row-1, leftSideCol-1)
		}
	} else if itemAboveRightSide == '[' {
		// . [ ]
		// [ ]
		return CanPushBoxUp(warehouse, row-1, leftSideCol+1)
	}

	return false
}

func PushBoxUp(warehouse [][]byte, row int, leftSideCol int) {
	itemAboveLeftSide := warehouse[row-1][leftSideCol]
	itemAboveRightSide := warehouse[row-1][leftSideCol+1]

	if itemAboveLeftSide == '[' {
		// [ ]
		// [ ]
		PushBoxUp(warehouse, row-1, leftSideCol)
	} else if itemAboveLeftSide == ']' {
		if itemAboveRightSide == '[' {
			// [ ] [ ]
			//   [ ]
			PushBoxUp(warehouse, row-1, leftSideCol-1)
			PushBoxUp(warehouse, row-1, leftSideCol+1)
		} else {
			// [ ] .
			//   [ ]
			PushBoxUp(warehouse, row-1, leftSideCol-1)
		}
	} else if itemAboveRightSide == '[' {
		// . [ ]
		// [ ]
		PushBoxUp(warehouse, row-1, leftSideCol+1)
	}

	warehouse[row-1][leftSideCol] = '['
	warehouse[row-1][leftSideCol+1] = ']'
	warehouse[row][leftSideCol] = '.'
	warehouse[row][leftSideCol+1] = '.'
}

func CanPushBoxDown(warehouse [][]byte, row int, leftSideCol int) bool {
	itemBelowLeftSide := warehouse[row+1][leftSideCol]
	itemBelowRightSide := warehouse[row+1][leftSideCol+1]

	// [ ]  or  [ ]  or  [ ]
	// # #      #          #
	if itemBelowLeftSide == '#' || itemBelowRightSide == '#' {
		return false
	}

	// . .
	// [ ]
	if itemBelowLeftSide == '.' && itemBelowRightSide == '.' {
		return true
	}

	if itemBelowLeftSide == '[' {
		// [ ]
		// [ ]
		return CanPushBoxDown(warehouse, row+1, leftSideCol)
	} else if itemBelowLeftSide == ']' {
		if itemBelowRightSide == '[' {
			//   [ ]
			// [ ] [ ]
			return CanPushBoxDown(warehouse, row+1, leftSideCol-1) && CanPushBoxDown(warehouse, row+1, leftSideCol+1)
		} else {
			//   [ ]
			// [ ] .
			return CanPushBoxDown(warehouse, row+1, leftSideCol-1)
		}
	} else if itemBelowRightSide == '[' {
		// [ ]
		// . [ ]
		return CanPushBoxDown(warehouse, row+1, leftSideCol+1)
	}

	return false
}

func PushBoxDown(warehouse [][]byte, row int, leftSideCol int) {
	itemBelowLeftSide := warehouse[row+1][leftSideCol]
	itemBelowRightSide := warehouse[row+1][leftSideCol+1]

	if itemBelowLeftSide == '[' {
		// [ ]
		// [ ]
		PushBoxDown(warehouse, row+1, leftSideCol)
	} else if itemBelowLeftSide == ']' {
		if itemBelowRightSide == '[' {
			//   [ ]
			// [ ] [ ]
			PushBoxDown(warehouse, row+1, leftSideCol-1)
			PushBoxDown(warehouse, row+1, leftSideCol+1)
		} else {
			//   [ ]
			// [ ] .
			PushBoxDown(warehouse, row+1, leftSideCol-1)
		}
	} else if itemBelowRightSide == '[' {
		// [ ]
		// . [ ]
		PushBoxDown(warehouse, row+1, leftSideCol+1)
	}

	warehouse[row+1][leftSideCol] = '['
	warehouse[row+1][leftSideCol+1] = ']'
	warehouse[row][leftSideCol] = '.'
	warehouse[row][leftSideCol+1] = '.'
}

func SumBoxGpsCoordinates(warehouse [][]byte) int {
	runningTotal := 0
	for row, items := range warehouse {
		for col, item := range items {
			if item == '[' {
				runningTotal += 100*(row-1) + col
			}
		}
	}

	return runningTotal
}
