package main

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/RomanKovalev007/webmap/internal/worker"
)



func main() {
	startURL := "https://simple.wikipedia.org/"
	maxWorkers := 5
	maxDepth := 2
	maxPages := 500
	timeout := 5 * time.Second

	var mutex sync.Mutex

	visited := make(map[string]bool)
	mutex.Lock()
	visited[startURL] = true
	mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	parsedURL, err := url.Parse(startURL)
	if err != nil {
		fmt.Printf("Invalid start URL: %v\n", err)
		return
	}
	baseDomain := parsedURL.Hostname()

	tasks := make(chan worker.CrawlTask, 1000)
	results := worker.WorkerPool(ctx, tasks, maxDepth, baseDomain, maxWorkers)


	tasks <- worker.CrawlTask{
		URL: startURL,
		Depth: 0,
	}
	fmt.Println("первая задача записана в таски")

	for task := range results{
		//fmt.Println("MAIN пришла из результатов ", task)			
		mutex.Lock()
		if len(visited) > maxPages{
			mutex.Unlock()
			cancel()
			close(tasks)
			fmt.Println("канал задач закрылся")
			break 
		}
		
		if task.Depth > maxDepth {
			mutex.Unlock()
			continue
		}

		if visited[task.URL]{
			mutex.Unlock()
			continue
		}
		visited[task.URL] = true
		mutex.Unlock()


		tasks <- task
		//fmt.Println("MAIN отправилась в таски ", task)
	}

	fmt.Println("Найдено страниц: ", len(visited))
}



	