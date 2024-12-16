package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
)

func parseInputFile(path string) [][]byte {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		panic(0)
	}
	defer file.Close()

	var maze [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) <= 0 {
			break
		}
		maze = append(maze, []byte(line))
	}

	return maze
}

func main() {
	var path string
	flag.StringVar(&path, "path", "", "The path to the input file")

	var debug bool
	flag.BoolVar(&debug, "debug", false, "Specify whether or not to produce debug output")

	flag.Parse()

	maze := parseInputFile(path)

	if debug {
		printMaze(maze)
	}

	reindeer := findReindeer(maze)
	end := findEnd(maze)

	shortestRoutes := findShortestRoutes(reindeer, end, maze, debug)
	fmt.Println("Lowest possible score found to be ", shortestRoutes[0].Score)
	fmt.Println("Tiles on shortest routes found to be ", countTilesOnAnyShortestRoute(maze, shortestRoutes))
}

func countTilesOnAnyShortestRoute(maze [][]byte, shortestRoutes []Route) int {
	traversedMaze := make([][]byte, len(maze))
	for row, items := range maze {
		itemsCopy := make([]byte, len(items))
		copy(itemsCopy, items)
		traversedMaze[row] = itemsCopy
	}

	for _, route := range shortestRoutes {
		for _, step := range route.Path {
			traversedMaze[step.Point.Row][step.Point.Col] = 'O'
		}
	}

	tileCount := 0
	for _, items := range traversedMaze {
		for _, item := range items {
			if item == 'O' {
				tileCount++
			}
		}
	}

	return tileCount
}

func findReindeer(maze [][]byte) DirectedPoint {
	for row, items := range maze {
		for col, item := range items {
			if item == 'S' {
				return DirectedPoint{
					Point: Point{
						Row: row,
						Col: col,
					},
					Direction: East,
				}
			}
		}
	}

	fmt.Println(fmt.Errorf("Could not find reindeer in maze"))
	panic(0)
}

func findEnd(maze [][]byte) Point {
	for row, items := range maze {
		for col, item := range items {
			if item == 'E' {
				return Point{
					Row: row,
					Col: col,
				}
			}
		}
	}

	fmt.Println(fmt.Errorf("Could not find end in maze"))
	panic(0)
}

func printMaze(maze [][]byte) {
	for _, items := range maze {
		fmt.Println(string(items[:]))
	}
	fmt.Print("\n")
}

func printRoute(maze [][]byte, route Route) {
	traversedMaze := make([][]byte, len(maze))
	for row, items := range maze {
		itemsCopy := make([]byte, len(items))
		copy(itemsCopy, items)
		traversedMaze[row] = itemsCopy
	}

	for _, step := range route.Path {
		traversedMaze[step.Point.Row][step.Point.Col] = byte(step.Direction)
	}

	fmt.Println("Score: ", route.Score)
	printMaze(traversedMaze)
}

type direction byte

const (
	North = '^'
	East  = '>'
	South = 'v'
	West  = '<'
)

type Point struct {
	Row int
	Col int
}

type DirectedPoint struct {
	Point     Point
	Direction direction
}

type Route struct {
	Reindeer DirectedPoint
	Path     []DirectedPoint
	Score    int
}

func findShortestRoutes(start DirectedPoint, end Point, maze [][]byte, debug bool) []Route {
	winningScore := math.MaxInt
	shortestRoutes := make([]Route, 0)

	queue := []Route{{
		Reindeer: start,
		Path:     []DirectedPoint{start},
		Score:    0,
	}}
	lowestScores := make(map[DirectedPoint]int)

	for len(queue) > 0 {
		route := queue[0]
		queue = queue[1:]

		if debug {
			printRoute(maze, route)
		}

		if route.Score > winningScore {
			continue
		}

		if route.Reindeer.Point == end {
			if route.Score < winningScore {
				winningScore = route.Score
				shortestRoutes = []Route{route}
			} else if route.Score == winningScore {
				shortestRoutes = append(shortestRoutes, route)
			}

			continue
		}

		for _, step := range getPossibleReindeerSteps(route.Reindeer, maze) {
			scoreAfterStep := route.Score + 1
			if step.Direction != route.Reindeer.Direction {
				scoreAfterStep += 1000
			}

			if lowestScore, visited := lowestScores[step]; visited {
				if scoreAfterStep > lowestScore {
					continue
				}
			}

			lowestScores[step] = scoreAfterStep

			pathAfterStep := make([]DirectedPoint, len(route.Path))
			copy(pathAfterStep, route.Path)
			pathAfterStep = append(pathAfterStep, step)
			queue = append(queue, Route{
				Reindeer: step,
				Path:     pathAfterStep,
				Score:    scoreAfterStep,
			})
		}
	}

	return shortestRoutes
}

func getPossibleReindeerSteps(reindeer DirectedPoint, maze [][]byte) []DirectedPoint {
	possibleReindeerSteps := make([]DirectedPoint, 0)

	north := getAdjacentPoint(reindeer.Point, North)
	east := getAdjacentPoint(reindeer.Point, East)
	south := getAdjacentPoint(reindeer.Point, South)
	west := getAdjacentPoint(reindeer.Point, West)

	if reindeer.Direction == North {
		if isOnPath(north, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{north, North})
		}
		if isOnPath(east, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{east, East})
		}
		if isOnPath(west, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{west, West})
		}
	}

	if reindeer.Direction == East {
		if isOnPath(north, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{north, North})
		}
		if isOnPath(east, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{east, East})
		}
		if isOnPath(south, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{south, South})
		}
	}

	if reindeer.Direction == South {
		if isOnPath(east, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{east, East})
		}
		if isOnPath(south, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{south, South})
		}
		if isOnPath(west, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{west, West})
		}
	}

	if reindeer.Direction == West {
		if isOnPath(north, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{north, North})
		}
		if isOnPath(south, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{south, South})
		}
		if isOnPath(west, maze) {
			possibleReindeerSteps = append(possibleReindeerSteps, DirectedPoint{west, West})
		}
	}

	return possibleReindeerSteps
}

func isOnPath(point Point, maze [][]byte) bool {
	item := maze[point.Row][point.Col]
	return item == '.' || item == 'E'
}

func getAdjacentPoint(point Point, dir direction) Point {
	if dir == North {
		return Point{Row: point.Row - 1, Col: point.Col}
	} else if dir == East {
		return Point{Row: point.Row, Col: point.Col + 1}
	} else if dir == South {
		return Point{Row: point.Row + 1, Col: point.Col}
	} else {
		return Point{Row: point.Row, Col: point.Col - 1}
	}
}
