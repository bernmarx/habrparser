package scraper

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bernmarx/habrparser/app/internal/page"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseOfNum        = 10
	linkClass        = "tm-article-snippet__title-link"
	authorClass      = ".tm-article-snippet__author"
	titleClass       = ".tm-article-snippet__title_h1"
	timeTag          = "time"
	commentClass     = ".tm-article-page-comments"
	ratingClass      = ".tm-votes-meter__value_appearance-article"
	articleTextClass = ".article-formatted-body"
)

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Scraper struct {
	HttpClient
}

func NewScraper() *Scraper {
	return &Scraper{HttpClient: http.DefaultClient}
}

func (s *Scraper) SetClient(client *http.Client) {
	s.HttpClient = client
}

// ScrapeLinks извлечение ссылок на статьи из ежедневного топа
func (s *Scraper) ScrapeLinks(url string, maxPages int) ([]string, error) {
	parsedLinks := make([]string, 0, maxPages)
	res, err := s.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("ScrapeLinks error: Status code is not OK")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	doc.Find("." + linkClass).Each(func(i int, s *goquery.Selection) {
		if i >= maxPages {
			return
		}

		link, exists := s.Attr("href")
		if !exists {
			log.Println("article link was selected but no link found")
			err = errors.New("article link was selected but no link found")
		}

		parsedLinks = append(parsedLinks, link)
	})

	return parsedLinks, err
}

// ScrapeArticle извлекает данные из статьи с хабра
func (s *Scraper) ScrapeArticle(pageURL string) (page.Page, error) {
	parsedPage := page.Page{}
	res, err := s.Get(pageURL)

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
		return parsedPage, err
	}

	id, err := strconv.ParseInt(pageURL[len(pageURL)-7:len(pageURL)-1], baseOfNum, 0)
	if err != nil {
		log.Fatalf("failed to parse id from %v", pageURL)
	}

	parsedPage.ID = int(id)

	parsedPage.Author, err = findAuthor(doc)
	if err != nil {
		return parsedPage, err
	}

	parsedPage.Title = findTitle(doc)

	parsedPage.Posted, err = findPosted(doc)
	if err != nil {
		return parsedPage, err
	}

	parsedPage.CommentCount, err = findCommentCount(doc)
	if err != nil {
		return parsedPage, err
	}

	parsedPage.Rating, err = findRating(doc)
	if err != nil {
		return parsedPage, err
	}

	parsedPage.Article = findArticleText(doc)

	return parsedPage, nil
}

func findAuthor(doc *goquery.Document) (string, error) {
	selection := doc.Find(authorClass)
	if selection == nil {
		return "", errors.New("could not find author")
	}

	return strings.TrimSpace(selection.Text()), nil
}
func findTitle(doc *goquery.Document) string {
	selection := doc.Find(titleClass)

	return strings.TrimSpace(selection.Text())
}

func findPosted(doc *goquery.Document) (string, error) {
	selection := doc.Find(timeTag)

	posted, exists := selection.Attr("datetime")
	if exists {
		return posted, nil
	}

	return "", errors.New("could not find posted")
}

func findCommentCount(doc *goquery.Document) (int, error) {
	selection := doc.Find(commentClass)

	n, err := getNumber(selection.Text())
	if err != nil {
		return 0, err
	}
	return n, nil
}

func findRating(doc *goquery.Document) (int, error) {
	selection := doc.Find(ratingClass).First()

	rating, err := strconv.ParseInt(selection.Text(), baseOfNum, 0)
	if err != nil {
		if selection.Text() == "" {
			return 0, nil
		}
		return 0, errors.New("failed to parse rating from article")
	}

	return int(rating), nil
}
func findArticleText(doc *goquery.Document) string {
	selection := doc.Find(articleTextClass)

	return selection.Text()
}

func getNumber(str string) (int, error) {
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
			return 0, nil
		}
		return 0, errors.New("failed to parse number")
	}

	return int(ans), nil
}
