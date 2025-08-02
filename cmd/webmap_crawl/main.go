package main

import (
	"fmt"
	"sync"

	"github.com/RomanKovalev007/webmap/internal/worker"
)



func main() {
	startURL := ""
	maxWorkers := 5

	tasks := make(chan string)
	results := make(chan string)
	visited := make(map[string]bool)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	for i := 1; i <= maxWorkers; i++ {
		wg.Add(1)
		go worker.Worker(tasks, results, visited, &wg, &mutex)
	}

	go func()  {
		tasks <- startURL
	}()

	go func()  {
		for link := range results{
			mutex.Lock()
			if !visited[link] {
				tasks <- link
			}
			mutex.Unlock()
		}
		close(tasks)
	}()

	wg.Wait()
	close(results)

	fmt.Println("Crawling finished!")
	fmt.Println("Visited pages:", len(visited))
	for url := range visited{
		fmt.Println("-", url)
	}
}