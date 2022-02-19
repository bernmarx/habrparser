package main

import (
	"database/sql"
	"log"

	"github.com/bernmarx/habrparser/internal/pkg/page"
	"github.com/bernmarx/habrparser/internal/pkg/scraper"

	_ "github.com/lib/pq"
)

const (
	dailyURL = "https://habr.com/ru/top/daily/"
	maxPages = 10
)

func worker(jobs <-chan string, results chan<- page.Page) {
	for j := range jobs {
		results <- scraper.ScrapeArticle(j)
	}
}

func main() {
	//подключение к бд
	DataBase, err := sql.Open("postgres",
		"host=db port=5432 user=username password=password dbname=habr_pages sslmode=disable")
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	err = DataBase.Ping()
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	log.Println("connected to database!")

	//парсинг ссылок
	parsedLinks := scraper.ScrapeLinks(dailyURL, maxPages, "tm-article-snippet__title-link")

	jobs := make(chan string, maxPages)
	results := make(chan page.Page, maxPages)

	log.Println(parsedLinks)

	//парсинг статей и их отправление
	for i := 0; i < maxPages; i++ {
		go worker(jobs, results)
		jobs <- "https://habr.com" + parsedLinks[i]
	}

	close(jobs)

	for i := 0; i < maxPages; i++ {
		r := <-results
		sqlstmt := `SELECT addpage($1)`
		sqlstmt2 := `SELECT addpagejson($1, $2)`
		j := r.GetJSON()

		_, err := DataBase.Exec(sqlstmt, j)
		if err != nil {
			log.Fatalf("failed to add page\n%v", err.Error())
		}

		_, err = DataBase.Exec(sqlstmt2, r.ID, j)

		//обработка ошибок
		if err != nil {
			log.Fatalf("failed to add page \n%v", err.Error())
		}
	}

	log.Println("finished parsing!")
}
