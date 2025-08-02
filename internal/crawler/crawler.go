package crawler

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func Crawler(url string) ([]string, error){
	
	resp, err := http.Get(url)
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

	var links []string
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a"{
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := attr.Val
					if strings.HasPrefix(link, url){
						links = append(links, link)
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