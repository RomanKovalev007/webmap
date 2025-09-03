package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/RomanKovalev007/webmap/internal/crawler"
)

// структура задачи
type CrawlTask struct {
    URL   string
    Depth int  // Текущая глубина (начинается с 0)
}

// функция основного обработчика
func Worker(
	ctx context.Context,
	mutex *sync.Mutex,
	tasks <-chan CrawlTask,
	results chan<- CrawlTask,
	visited map[string]bool,
	wg *sync.WaitGroup,
	maxPages int,
	baseDomain string,
	num int){

	defer wg.Done()

	for {
		// обработка контекста
		select {
		case <-ctx.Done():
			fmt.Println("Воркер номер ", num, "завершил свою работу по контексту")
			return
		// обработка задач поступающих из канала задач
		case task, ok := <-tasks:
			// если нет задач, функция перестает работать
			if !ok {
				fmt.Println("нет входящих задач")
				return
			}

			fmt.Println("WORKER обработка страницы ", task, "воркер ", num)
			
			// парсинг страницы с помощью функции краулера
			links, err := crawler.Crawler(task.URL, baseDomain)
			if err != nil {
				fmt.Println("Ошибка краулинга: ", err)
				continue
			}

			fmt.Println(len(links))
			
			// отправка найденных при парсинге ссылок в канал результатов
			for _, link := range links {
				select {
				case results <- CrawlTask{URL: link, Depth: task.Depth + 1}:
					mutex.Lock()
					if len(visited) >= maxPages{
						mutex.Unlock()
						fmt.Println("Воркер номер ", num, "завершил свою работу так как список посещенных страниц превысил свою длину")
						return
					}
					mutex.Unlock()
				case <-ctx.Done():
					fmt.Println("Воркер номер ", num, "завершил свою работу по контексту")
					return
				}
				
			} 
	}
	}
}

// пул запуска воркеров
func WorkerPool(
	ctx context.Context,
	mutex *sync.Mutex,
	tasks chan CrawlTask,
	visited map[string]bool,
	maxPages int,
	baseDomain string,
	WorkerCount int) chan CrawlTask {

	wg := sync.WaitGroup{}

	results := make(chan CrawlTask)

	go func ()  {
		for i := 1; i <= WorkerCount; i++{
			wg.Add(1)
			fmt.Println("воркер запущен ", i)
			go Worker(ctx, mutex, tasks, results, visited, &wg, maxPages, baseDomain, i)
		} 
		wg.Wait()
		fmt.Println("канал результатов закрылся")
		close(results)

	}()

	return results
}