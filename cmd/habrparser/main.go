package main

import (
	"log"
	"time"

	"github.com/bernmarx/habrparser/internal/scraper"
	"github.com/bernmarx/habrparser/internal/storage"
	"github.com/getsentry/sentry-go"

	_ "github.com/lib/pq"
)

const (
	dailyURL = "https://habr.com/ru/top/daily/"
	maxPages = 10
	workers  = 8
)

func worker(s *scraper.Scraper, jobs <-chan string, results chan<- scraper.Page) {
	for j := range jobs {
		r, err := s.ScrapeArticle(j)
		if err != nil {
			sentry.CaptureException(err)
			log.Println("worker encountered error\n", err.Error())
		}

		results <- r
	}
}

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "http://61eac3e4cf2d47fe966353777d026ffd@127.0.0.1:9000/1",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage("connected to sentry!")
	log.Println("connected to sentry")

	scr := scraper.NewScraper()

	//парсинг ссылок
	parsedLinks, err := scr.ScrapeLinks(dailyURL, maxPages)
	if err != nil {
		sentry.CaptureException(err)
	}

	jobs := make(chan string, maxPages)
	results := make(chan scraper.Page, maxPages)

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
		sentry.CaptureException(err)
		log.Fatalln("could not connect to database\n", err.Error())
	}

	for i := 0; i < maxPages; i++ {
		r := <-results

		s.AddPageData(&r)
	}
}
