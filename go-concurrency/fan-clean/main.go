package main

import (
	"fmt"
	"sync"
)

var numbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

const numWorkers = 3

// worker goroutine
func worker(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for number := range in {
			out <- number * number
		}
		close(out)
	}()

	return out
}

// fan in - channel merging goroutine
func fanIn(in []<-chan int) <-chan int {

	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(in))
	for _, ch := range in {
		go func() {
			for number := range ch {
				out <- number
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {

	// feeder
	input := make(chan int)
	go func() {
		for _, i := range numbers {
			input <- i
		}
		close(input)
	}()

	// pool of "clean" workers
	var workerChannels []<-chan int
	for range numWorkers {
		workerChannels = append(workerChannels, worker(input))
	}

	// collect results
	for result := range fanIn(workerChannels) {
		fmt.Println(result)
	}
}
