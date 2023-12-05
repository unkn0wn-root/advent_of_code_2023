package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type RequirementRange struct {
	Source      int
	Destination int
	Length      int
	Multiplier  int
}

type SeedRequirement struct {
	FromDest     string
	ToDest       string
	Requirements []RequirementRange
}

// nextId returns the next Id based on the given seed within the range.
func (r *RequirementRange) nextId(seed int) int {
	if seed >= r.Source && seed < (r.Source+r.Length) {
		return r.Destination + seed - r.Source
	}

	return -1
}

// returns the next requirement ID based on the given seed for the seed requirement.
func (s *SeedRequirement) getNextReqId(seed int) int {
	for _, req := range s.Requirements {
		nextID := req.nextId(seed)
		if nextID != -1 {
			return nextID
		}
	}

	return seed
}

// getSeeds extracts integers from the given byte slice.
func getSeeds(data []byte) (out []int) {
	reg := regexp.MustCompile(`\d+`)
	for _, rawNum := range reg.FindAll(data, -1) {
		n, _ := strconv.Atoi(string(rawNum))
		out = append(out, n)
	}

	return
}

// extracts "from" and "to" destinations from the input data.
func getMap(data [][]byte, startIndex int) (fromDest string, toDest string) {
	mapRegexp := regexp.MustCompile(`(\w+)-to-(\w+).*`)
	matched := mapRegexp.FindSubmatch(data[startIndex])
	return string(matched[1]), string(matched[2])
}

// extracts RequirementRange from the input data.
func getRangeRequirements(data [][]byte, startIndex int) RequirementRange {
	rangeRegexp := regexp.MustCompile(`(\d+)`)
	matched := rangeRegexp.FindAllSubmatch(data[startIndex], -1)

	destination, source, length := matched[0][0], matched[1][0], matched[2][0]
	destInt, _ := strconv.Atoi(string(destination))
	sourceInt, _ := strconv.Atoi(string(source))
	lengthInt, _ := strconv.Atoi(string(length))

	return RequirementRange{
		Source:      sourceInt,
		Destination: destInt,
		Length:      lengthInt,
	}
}

// extracts SeedRequirements from the input data.
func getSeedRequirements(data [][]byte) (out []SeedRequirement) {
	startIndex := 2

	for startIndex < len(data) {
		fromDest, toDest := getMap(data, startIndex)
		seedReq := SeedRequirement{FromDest: fromDest, ToDest: toDest}

		startIndex++
		for startIndex < len(data) && len(data[startIndex]) > 0 {
			seedReq.Requirements = append(seedReq.Requirements, getRangeRequirements(data, startIndex))
			startIndex++
		}

		out = append(out, seedReq)
		startIndex++
	}

	return
}

func partOne(seeds []int, seedRequirements []SeedRequirement) (int, time.Duration) {
	lowestLocation := math.Inf(1)

	startTime := time.Now()

	for _, seed := range seeds {
		for _, seedReq := range seedRequirements {
			seed = seedReq.getNextReqId(seed)
		}

		if float64(seed) < lowestLocation {
			lowestLocation = float64(seed)
		}
	}

	finishTime := time.Since(startTime)

	return int(lowestLocation), finishTime
}

// calculates the lowest location using goroutines to speed up thing a bit. Not ideal though.
func partTwo(seeds []int, seedRequirements []SeedRequirement) (int, time.Duration) {
	var lowestLocationMutex sync.Mutex
	lowestLocation := math.Inf(1)

	var wg sync.WaitGroup

	startTime := time.Now()

	for seedIndex := 0; seedIndex < len(seeds); seedIndex += 2 {
		wg.Add(1)

		go func(seedIndex int) {
			defer wg.Done()
			startSeed, seedLength := seeds[seedIndex], seeds[seedIndex+1]
			offsets := make([]int, seedLength)

			for i := 0; i < seedLength; i++ {
				offsets[i] = i
			}

			for _, seedReq := range seedRequirements {
				for i := 0; i < seedLength; i++ {
					offsets[i] = seedReq.getNextReqId(startSeed+offsets[i]) - startSeed
				}
			}

			for _, offset := range offsets {
				result := startSeed + offset
				lowestLocationMutex.Lock()

				if float64(result) < lowestLocation {
					lowestLocation = float64(result)
				}

				lowestLocationMutex.Unlock()
			}
		}(seedIndex)
	}

	wg.Wait()

	finishTime := time.Since(startTime)

	return int(lowestLocation), finishTime
}

func getLocalInputFile(inputPath string) (in [][]byte) {
	file, errFile := os.Open(inputPath)
	if errFile != nil {
		panic(errFile)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		in = append(in, []byte(row))
	}

	return
}

func main() {
	data := getLocalInputFile("input.txt")
	seeds := getSeeds(data[0])
	seedRequirements := getSeedRequirements(data)

	partOneCalculation, timeTaken := partOne(seeds, seedRequirements)
	fmt.Println("Part One:", partOneCalculation, "Time taken:", timeTaken)

	fmt.Println("Wait for part two to finish calculating seeds...")

	partTwoCalculation, timeTaken := partTwo(seeds, seedRequirements)
	fmt.Println("Part Two:", partTwoCalculation, "Time taken:", timeTaken.Truncate(time.Second))
}
