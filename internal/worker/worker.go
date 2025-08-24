package worker

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/RomanKovalev007/webmap/internal/crawler"
)

// структура задачи
type CrawlTask struct {
    URL   string
    Depth int  // Текущая глубина (начинается с 0)
}

// функция основного обработчика
func Worker(ctx context.Context, tasks <-chan CrawlTask, results chan<- CrawlTask, wg *sync.WaitGroup, maxDepth int, baseDomain string){
	defer wg.Done()

	for {
		// обработка контекста
		select {
		case <-ctx.Done():
			return
		case task, ok := <-tasks:
			if !ok {
				fmt.Println("нет входящих задач")
				return
			}
			fmt.Println("обработка страницы ", task)

			links, err := crawler.Crawler(task.URL)
			if err != nil {
				fmt.Println("Ошибка краулинга: ", err)
				return
			}

			for _, link := range links {
				parsedURL, err := url.Parse(link)
				if err != nil {
					fmt.Println("Ошибка распарсивания ссылки в воркере")
					continue
				}

				if strings.HasSuffix(parsedURL.Hostname(), baseDomain) {
					select {
					case results <- CrawlTask{URL: link, Depth: task.Depth + 1}:
					case <-ctx.Done():
						return
					}
				}
			}
	}
	}
}