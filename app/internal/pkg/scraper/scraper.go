package scraper

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bernmarx/habrparser/app/internal/pkg/page"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseOfNum = 10
)

// ScrapeLinks извлечение ссылок на статьи из ежедневного топа
func ScrapeLinks(url string, maxPages int, class string) []string {
	parsedLinks := make([]string, 0, maxPages)
	res, err := http.Get(url)

	if err != nil {
		log.Fatalf("%v\n", err.Error())
	}
	//processErr(err)

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	doc.Find("." + class).Each(func(i int, s *goquery.Selection) {
		if i >= maxPages {
			return
		}

		link, exists := s.Attr("href")
		if !exists {
			log.Fatal("article link was selected but no link was found")
		}

		parsedLinks = append(parsedLinks, link)
	})

	return parsedLinks
}

func ScrapeArticle(pageURL string) page.Page {
	parsedPage := page.Page{}
	res, err := http.Get(pageURL)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("\nstatus code error: %d %s while connecting to %v", res.StatusCode, res.Status, pageURL)
	}

	if res.StatusCode == http.StatusOK {
		log.Println("\nconnected to ", pageURL)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatalf("%v\n", err.Error())
	}

	id, err := strconv.ParseInt(pageURL[len(pageURL)-7:len(pageURL)-1], baseOfNum, 0)

	if err != nil {
		log.Fatalf("failed to parse id from %v", pageURL)
	}

	parsedPage.ID = int(id)
	parsedPage.Author = findAuthor(doc)
	parsedPage.Title = findTitle(doc)
	parsedPage.Posted = findPosted(doc)
	parsedPage.CommentCount = findCommentCount(doc)
	parsedPage.Rating = findRating(doc)
	parsedPage.Article = findArticleText(doc)

	return parsedPage
}

func findAuthor(doc *goquery.Document) string {
	selection := doc.Find(".tm-article-snippet__author")
	if selection == nil {
		log.Fatalln("no author found")
	}

	return strings.TrimSpace(selection.Text())
}
func findTitle(doc *goquery.Document) string {
	selection := doc.Find(".tm-article-snippet__title_h1")

	return strings.TrimSpace(selection.Text())
}

func findPosted(doc *goquery.Document) string {
	selection := doc.Find("time")

	posted, exists := selection.Attr("datetime")
	if exists {
		return posted
	}

	log.Fatal("no datetime found")

	return ""
}

func findCommentCount(doc *goquery.Document) int {
	selection := doc.Find(".tm-article-page-comments")

	return getNumber(selection.Text())
}

func findRating(doc *goquery.Document) int {
	selection := doc.Find(".tm-votes-meter__value_appearance-article").First()

	rating, err := strconv.ParseInt(selection.Text(), baseOfNum, 0)
	if err != nil {
		if selection.Text() == "" {
			return 0
		}
		log.Fatalf("failed to parse rating from %v", selection.Text())
	}

	return int(rating)
}
func findArticleText(doc *goquery.Document) string {
	selection := doc.Find(".article-formatted-body")

	return selection.Text()
}

func getNumber(str string) int {
	num := ""

	for len(str) > 0 {
		r, size := utf8.DecodeRuneInString(str)
		if r >= '0' && r <= '9' {
			num += string(r)
		}

		str = str[size:]
	}

	ans, err := strconv.ParseInt(num, baseOfNum, 0)
	if err != nil {
		if num == "" {
			return 0
		}
		log.Fatalf("failed to parse comment count from %v", str)
	}

	return int(ans)
}
