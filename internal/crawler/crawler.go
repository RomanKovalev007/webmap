package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// функция поиска подходящих ссылок на страницы
// на вход подается ссылка на обрабатываемую страницу
// на выходе получаем слайс подходящих ссылок
func Crawler(pageUrl string, baseDomain string) ([]string, error){
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}
	req, err := http.NewRequest("GET", pageUrl, nil)
	if err != nil{
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyCrawler/1.0)")

	resp, err := client.Do(req)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	html_doc, err := html.Parse(resp.Body)
	if err != nil{
		return nil, err
	}

	baseURL, err := url.Parse(pageUrl)
	if err != nil {
		return nil, fmt.Errorf("URL parse error: %v", err)
	}


	var links []string
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a"{
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := attr.Val
					if link == "" || strings.HasPrefix(link, "#") || strings.HasPrefix(link, "javascript:") {
						continue
					}
					parsedLink, err := url.Parse(link)
					if err != nil {
						continue
					}
					absLink := baseURL.ResolveReference(parsedLink).String()

					if strings.HasPrefix(absLink, "http"){
						links = append(links, absLink)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling{
			f(c)
		}
	}
	f(html_doc)

	return links, nil
}