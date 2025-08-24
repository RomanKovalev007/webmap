package main

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RomanKovalev007/webmap/internal/worker"
)



func main() {
	startURL := "https://simple.wikipedia.org/"
	maxWorkers := 5
	maxDepth := 2
	maxPages := 500
	timeout := 1 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	parsedURL, err := url.Parse(startURL)
	if err != nil {
		fmt.Printf("Invalid start URL: %v\n", err)
		return
	}

	baseDomain := parsedURL.Hostname()

	tasks := make(chan worker.CrawlTask, 100)
	results := make(chan worker.CrawlTask, 100)

	visited := make(map[string]bool)

	var mutex sync.Mutex
	var wg sync.WaitGroup

	var activeTasks int32 = 1

	// Отправка начальной задачи
	tasks <- worker.CrawlTask{
		URL: startURL,
		Depth: 0,
	}
	fmt.Println("первая задача записана в таски")

	// Запуск воркеров

	for i := 1; i <= maxWorkers; i++ {
		wg.Add(1)
		fmt.Println("воркер запущен ", i)
			time.Sleep(100 * time.Microsecond)
		go worker.Worker(ctx, tasks, results, &wg, maxDepth, baseDomain)
	}

	// Обработка результатов
    go func() {
        for {
			select {
			case <-ctx.Done():
				return
			case task, ok := <-results:
				if !ok {
					return
				}

            	mutex.Lock()
				shouldSend := !visited[task.URL] && len(visited) <= maxPages
				mutex.Unlock()

				if shouldSend {
					mutex.Lock()
					visited[task.URL] = true
					mutex.Unlock()
					
					fmt.Printf("Found %s (Depth %d)\n", task.URL, task.Depth)
					atomic.AddInt32(&activeTasks, 1)

					select {
					case <-ctx.Done():
					return
					case tasks <- task:
}
            }
            
        	}
		}
    }()

	go func() {
		for range tasks {
			if atomic.AddInt32(&activeTasks, -1) == 0 {
				close(tasks)
			}
		}
	}()
	
	
	wg.Wait()
	close(results)
	
	fmt.Println("Crawling finished!")
	fmt.Println("Visited pages:", len(visited))

}