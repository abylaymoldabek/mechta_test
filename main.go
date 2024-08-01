package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Numbers struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	numGoroutines := flag.Int("goroutines", 1, "Number of goroutines for parallel processing")
	flag.Parse()

	file, err := os.Open("data.json")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var numbers []Numbers
	err = json.Unmarshal(byteValue, &numbers)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	sum := calculateSum(numbers, *numGoroutines)

	fmt.Printf("Total sum: %d\n", sum)
}

func calculateSum(numbers []Numbers, numGoroutines int) int {
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	sumChan := make(chan int, numGoroutines)

	chunkSize := len(numbers) / numGoroutines
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if i == numGoroutines-1 {
			end = len(numbers)
		}

		go func(start, end int) {
			defer wg.Done()
			localSum := 0
			for j := start; j < end; j++ {
				localSum += numbers[j].A + numbers[j].B
			}
			sumChan <- localSum
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(sumChan)
	}()

	totalSum := 0
	for sum := range sumChan {
		totalSum += sum
	}

	return totalSum
}
