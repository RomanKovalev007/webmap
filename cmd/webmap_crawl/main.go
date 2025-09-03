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
	startURL := "https://mai.ru/"
	maxPages := 100
	maxDepth := 2
	timeout := 100 * time.Second
	maxWorkers := 5

	var mutex sync.Mutex

	visited := make(map[string]bool)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	parsedURL, err := url.Parse(startURL)
	if err != nil {
		fmt.Printf("Invalid start URL: %v\n", err)
		return
	}
	baseDomain := parsedURL.Hostname()

	tasks := make(chan worker.CrawlTask, maxPages)
	results := worker.WorkerPool(ctx, &mutex, tasks, visited, maxPages, baseDomain, maxWorkers)

	//записываем первую задачу в канал задач
	tasks <- worker.CrawlTask{URL: startURL, Depth: 0}
	mutex.Lock()
	visited[startURL] = true
	mutex.Unlock()

	//читаем из канала результатов и подходящие записываем в результирующую мапу
	for task := range results{
			mutex.Lock()
		
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
	}

	close(tasks)
	fmt.Println("канал задач закрылся")

			
	fmt.Println("Краулер завершил свою работу!")
	fmt.Println("Найдено страниц: ", len(visited))
}



	