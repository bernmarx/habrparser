package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func RequestWithProxy(method string, url string, body io.Reader, proxyURL *url.URL) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		log.Fatalf("got error while making request: %v", err.Error())
	}

	client := MakeProxyClient(proxyURL)
	defer client.CloseIdleConnections()

	return client.Do(req)
}

func MakeProxyClient(proxyURL *url.URL) http.Client {
	tr := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := http.Client{
		Transport: tr,
	}

	return client
}
