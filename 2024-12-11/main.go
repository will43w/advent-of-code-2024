package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var blinks int
	flag.IntVar(&blinks, "blinks", 1, "The number of times to blink")

	var input string
	flag.StringVar(&input, "input", "125 17", "The input array of stones, as a single string")

	flag.Parse()
	fmt.Println("Blinking ", blinks, " time(s) for the stone array ", input)

	stones := strings.Split(input, " ")
	stoneCounts := make(map[string]uint64)

	for _, stone := range stones {
		AddOrIncrementStoneCount(stoneCounts, stone, 1)
	}

	for blink := 0; blink < blinks; blink++ {
		stoneCounts = Blink(stoneCounts)
	}

	fmt.Println("Total number of resulting stones: ", GetTotalStoneCount(stoneCounts))
}

func AddOrIncrementStoneCount(stoneCounts map[string]uint64, stone string, increment uint64) {
	_, ok := stoneCounts[stone]
	if !ok {
		stoneCounts[stone] = increment
	} else {
		stoneCounts[stone] += increment
	}
}

func Blink(stoneCounts map[string]uint64) map[string]uint64 {
	newStoneCounts := make(map[string]uint64)
	for stone, count := range stoneCounts {
		if stone == "0" {
			AddOrIncrementStoneCount(newStoneCounts, "1", count)
		} else if len(stone)%2 == 0 {
			leftStone := strings.TrimLeft(stone[0:len(stone)/2], "0")
			if leftStone == "" {
				leftStone = "0"
			}

			rightStone := strings.TrimLeft(stone[len(stone)/2:], "0")
			if rightStone == "" {
				rightStone = "0"
			}

			AddOrIncrementStoneCount(newStoneCounts, leftStone, count)
			AddOrIncrementStoneCount(newStoneCounts, rightStone, count)
		} else {
			stoneuint64, err := strconv.Atoi(stone)
			if err != nil {
				fmt.Println(err)
				panic(0)
			}

			newStone := strings.TrimLeft(strconv.Itoa(stoneuint64*2024), "0")
			AddOrIncrementStoneCount(newStoneCounts, newStone, count)
		}
	}

	return newStoneCounts
}

func GetTotalStoneCount(stoneCounts map[string]uint64) uint64 {
	var totalStoneCount uint64 = 0
	for _, count := range stoneCounts {
		totalStoneCount += count
	}

	return totalStoneCount
}
