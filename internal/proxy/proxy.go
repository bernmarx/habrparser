package proxy

import (
	"net/http"
	"net/url"

	"github.com/bernmarx/habrparser/internal/scraper"
)

type ProxyClient struct {
	scraper.HttpClient
}

func NewProxyClient(proxy *url.URL) *ProxyClient {
	tr := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	client := http.Client{
		Transport: tr,
	}

	return &ProxyClient{HttpClient: &client}
}
