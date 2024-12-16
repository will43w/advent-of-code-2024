package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func ParseInputFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		panic(0)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func main() {
	var path string
	flag.StringVar(&path, "path", "", "The path to the input file")

	flag.Parse()

	fmt.Println(path)

	garden := ParseInputFile(path)
	fencingCalculator := NewFencingCalculator(garden)
	fencingCalculator.ExamineGarden()
	fmt.Println("Total fencing price: ", fencingCalculator.CalculateTotalFencingPrice())
}

type GardenPlot struct {
	Plant     rune
	Enclosure int

	area      int
	perimiter int
}

func (gp *GardenPlot) CalculateFencingPrice() int {
	return gp.area * gp.perimiter
}

type FencingCalculator struct {
	garden      []string
	inspected   [][]bool
	gardenPlots []*GardenPlot
}

func NewFencingCalculator(garden []string) *FencingCalculator {
	numRows := len(garden)
	numColumns := len(garden[0])

	inspected := [][]bool{}

	for row := 0; row < numRows; row++ {
		inspected = append(inspected, make([]bool, numColumns))
	}

	gardenPlots := make([]*GardenPlot, 0)

	return &FencingCalculator{
		garden:      garden,
		inspected:   inspected,
		gardenPlots: gardenPlots,
	}
}

func (fc *FencingCalculator) ExamineGarden() {
	numRows := len(fc.garden)
	numColumns := len(fc.garden[0])

	inspected := [][]bool{}

	for row := 0; row < numRows; row++ {
		inspected = append(inspected, make([]bool, numColumns))
	}

	enclosureId := 0
	for row := 0; row < numRows; row++ {
		for column := 0; column < numColumns; column++ {
			if !fc.inspected[row][column] {
				enclosureId++
				plant := rune(fc.garden[row][column])
				plot := &GardenPlot{
					Plant:     plant,
					Enclosure: enclosureId,
				}
				fc.gardenPlots = append(fc.gardenPlots, plot)
				fc.FloodFill(row, column, plant, plot)
			}
		}
	}
}

func (fc *FencingCalculator) FloodFill(row int, column int, plant rune, plot *GardenPlot) {
	if fc.inspected[row][column] {
		return
	}

	if fc.garden[row][column] != byte(plant) {
		return
	}

	fc.inspected[row][column] = true

	plot.area += 1

	perimiterContribution := 0
	if row == 0 {
		perimiterContribution += 1
	} else {
		plantAbove := fc.garden[row-1][column]
		if plantAbove != byte(plant) {
			perimiterContribution += 1
		} else {
			fc.FloodFill(row-1, column, plant, plot)
		}
	}

	if column+1 == len(fc.garden) {
		perimiterContribution += 1
	} else {
		plantToRight := fc.garden[row][column+1]
		if plantToRight != byte(plant) {
			perimiterContribution += 1
		} else {
			fc.FloodFill(row, column+1, plant, plot)
		}
	}

	if row+1 == len(fc.garden[0]) {
		perimiterContribution += 1
	} else {
		plantBelow := fc.garden[row+1][column]
		if plantBelow != byte(plant) {
			perimiterContribution += 1
		} else {
			fc.FloodFill(row+1, column, plant, plot)
		}
	}

	if column == 0 {
		perimiterContribution += 1
	} else {
		plantToLeft := fc.garden[row][column-1]
		if plantToLeft != byte(plant) {
			perimiterContribution += 1
		} else {
			fc.FloodFill(row, column-1, plant, plot)
		}
	}

	plot.perimiter += perimiterContribution
}

func (fc *FencingCalculator) CalculateTotalFencingPrice() int {
	totalPrice := 0
	for _, plot := range fc.gardenPlots {
		totalPrice += plot.CalculateFencingPrice()
	}

	return totalPrice
}
