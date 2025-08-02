package worker

import (
	"fmt"
	"sync"

	"github.com/RomanKovalev007/webmap/internal/crawler"
)

func Worker(tasks <-chan string, results chan<- string, visited map[string]bool, wg *sync.WaitGroup, mutex *sync.Mutex){
	defer wg.Done()

	for url := range tasks {
		mutex.Lock()
		if visited[url]{
			mutex.Unlock()
			continue
		}
		visited[url] = true
		mutex.Unlock()

		links, err := crawler.Crawler(url)
		if err != nil {
			fmt.Printf("Error crawling %s: %v\n", url, err)
		}

		for _, link := range links {
			results <- link
		}
	}
}