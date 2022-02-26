package main

import (
	"log"

	"github.com/bernmarx/habrparser/app/internal/page"
	"github.com/bernmarx/habrparser/app/internal/scraper"
	"github.com/bernmarx/habrparser/app/internal/storage"

	_ "github.com/lib/pq"
)

const (
	dailyURL = "https://habr.com/ru/top/daily/"
	maxPages = 10
	workers  = 8
)

func worker(s *scraper.Scraper, jobs <-chan string, results chan<- page.Page) {
	for j := range jobs {
		r, err := s.ScrapeArticle(j)
		if err != nil {
			log.Println("worker encountered error\n", err.Error())
		}

		results <- r
	}
}

func main() {
	scr := scraper.NewScraper()

	//парсинг ссылок
	parsedLinks, err := scr.ScrapeLinks(dailyURL, maxPages)
	if err != nil {
		log.Println("encountered error while parsing links\n", err.Error())
	}

	jobs := make(chan string, maxPages)
	results := make(chan page.Page, maxPages)

	log.Println(parsedLinks)

	for i := 0; i < workers; i++ {
		tempScr := scr
		go worker(tempScr, jobs, results)
	}

	for i := 0; i < maxPages; i++ {
		jobs <- "https://habr.com" + parsedLinks[i]
	}

	close(jobs)

	s, err := storage.NewStorage()
	if err != nil {
		log.Fatalln("could not connect to database\n", err.Error())
	}

	for i := 0; i < maxPages; i++ {
		r := <-results

		s.AddPageData(&r)
	}

	log.Println("finished parsing!")
}
