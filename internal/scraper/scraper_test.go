package scraper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetClient(t *testing.T) {
	s := Scraper{}
	c := http.DefaultClient
	s.SetClient(c)
	if s.HttpClient != c {
		t.Error("httpClient should be same", c, s.HttpClient)
	}
}

func TestScrapeLinks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw,
				`<html>
				<body>
					<a href="test" class="tm-article-snippet__title-link"></a>
				</body>
				</html>`)
		},
	))
	defer ts.Close()

	s := NewScraper()

	str, err := s.ScrapeLinks(ts.URL, 1)
	if err != nil {
		t.Error("failed to scrape article\n", err.Error())
	}

	if str[0] != "test" {
		t.Error("found article link does not match expected \"test\"\n", str[0])
	}

	ts = httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw,
				`<html>
				<body>
					<a href="test" class="tm-article-snippet__title-link"></a>
				</body>
				</html>`)
		},
	))
	defer ts.Close()

	s = NewScraper()

	str, _ = s.ScrapeLinks(ts.URL, 0)
	if len(str) > 0 {
		t.Error("len(str) is supposed to be 0\n", str)
	}

	ts = httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw,
				`some garbage`)
		},
	))
	defer ts.Close()

	str, _ = s.ScrapeLinks(ts.URL, 1)
	if len(str) > 0 {
		t.Error("should not get any links")
	}

	_, err = s.ScrapeLinks("garbage", 1)
	if err == nil {
		t.Error("should get error because url is invalid")
	}

	ts = httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(400)
			fmt.Fprintf(rw,
				`some garbage`)
		},
	))
	defer ts.Close()

	_, err = s.ScrapeLinks(ts.URL, 1)
	if err == nil {
		t.Error("should get error because status code is not OK")
	}

	ts = httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw,
				`<html>
				<body>
					<a class="tm-article-snippet__title-link"></a>
				</body>
				</html>`)
		},
	))
	defer ts.Close()

	_, err = s.ScrapeLinks(ts.URL, 1)
	if err == nil {
		t.Error("should get error because there is no link")
	}
}

func TestScrapeArticle(t *testing.T) {
	// this first test will IGNORE absence of id since it is calculated from URL
	// and it is hard to replicate something similar to habr's url here
	ts := httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw,
				`<html>
				<body>
					<p class="tm-article-snippet__author">author</p>
					<h1 class="tm-article-snippet__title_h1">title</h1>
					<time datetime="time"></time>
					<p class="tm-article-page-comments">1</p>
					<p class="tm-votes-meter__value_appearance-article">2</p>
					<p class="article-formatted-body">text</p>
				</body>
				</html>
				`)
		},
	))
	defer ts.Close()

	s := NewScraper()
	expected := Page{
		Author:       "author",
		Title:        "title",
		Posted:       "time",
		Article:      "text",
		CommentCount: 1,
		Rating:       2,
	}
	got, err := s.ScrapeArticle(ts.URL)
	if err != nil {
		t.Error("got error while scraping article\n", err.Error())
	}

	if got != expected {
		t.Error("got is different than expected\n", got, "\n", expected)
	}

	ts = httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw,
				`garbage`)
		},
	))
	defer ts.Close()

	_, err = s.ScrapeArticle(ts.URL)
	if err == nil {
		t.Error("should get error since test server's content is garbage")
	}

	ts = httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(rw,
				`garbage`)
		},
	))
	defer ts.Close()

	_, err = s.ScrapeArticle(ts.URL)
	if err == nil {
		t.Error("should get error status code is not OK")
	}

	_, err = s.ScrapeArticle("123")
	if err == nil {
		t.Error("should get error url is invalid")
	}

	ts = httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw,
				`<html>
				<body>
					<h1 class="tm-article-snippet__title_h1">title</h1>
					<time datetime="time"></time>
					<p class="tm-article-page-comments">1</p>
					<p class="tm-votes-meter__value_appearance-article">2</p>
					<p class="article-formatted-body">text</p>
				</body>
				</html>
				`)
		},
	))
	defer ts.Close()

	_, err = s.ScrapeArticle(ts.URL)
	if err == nil {
		t.Error("should get error author is missing")
	}

	ts = httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw,
				`<html>
				<body>
					<p class="tm-article-snippet__author">author</p>
					<h1 class="tm-article-snippet__title_h1">title</h1>
					<p class="tm-article-page-comments">1</p>
					<p class="tm-votes-meter__value_appearance-article">2</p>
					<p class="article-formatted-body">text</p>
				</body>
				</html>
				`)
		},
	))
	defer ts.Close()

	_, err = s.ScrapeArticle(ts.URL)
	if err == nil {
		t.Error("should get error posted is missing")
	}

	ts = httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw,
				`<html>
				<body>
					<p class="tm-article-snippet__author">author</p>
					<h1 class="tm-article-snippet__title_h1">title</h1>
					<time datetime="time"></time>
					<p class="tm-article-page-comments">asd</p>
					<p class="tm-votes-meter__value_appearance-article">2</p>
					<p class="article-formatted-body">text</p>
				</body>
				</html>
				`)
		},
	))
	defer ts.Close()

	p, _ := s.ScrapeArticle(ts.URL)
	if p.CommentCount != 0 {
		fmt.Println("comment count: ", p.CommentCount)
		t.Error("comment count should be 0")
	}

	ts = httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw,
				`<html>
				<body>
					<p class="tm-article-snippet__author">author</p>
					<h1 class="tm-article-snippet__title_h1">title</h1>
					<time datetime="time"></time>
					<p class="tm-article-page-comments">1</p>
					<p class="article-formatted-body">text</p>
				</body>
				</html>
				`)
		},
	))
	defer ts.Close()

	_, err = s.ScrapeArticle(ts.URL)
	if err == nil {
		t.Error("should get error rating is missing")
	}
}
