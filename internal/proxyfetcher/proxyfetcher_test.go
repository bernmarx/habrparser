package proxyfetcher

import (
	"net/url"
	"path/filepath"
	"testing"
)

func TestNewProxyListFromFile(t *testing.T) {
	file, _ := filepath.Abs("./../../test/valid_proxies.txt")
	_, err := NewProxyListFromFile(file)
	if err != nil {
		t.Error("got error while executing NewProxyListFromFile")
	}

	file, _ = filepath.Abs("./../../test/invalid_proxies.txt")
	p, err := NewProxyListFromFile(file)
	if err == nil {
		t.Error("should get error", p)
	}

	_, err = NewProxyListFromFile("123")
	if err == nil {
		t.Error("should get error")
	}
}

func TestFetchRandomProxy(t *testing.T) {
	file, _ := filepath.Abs("./../../test/one_proxy.txt")
	pl, err := NewProxyListFromFile(file)
	if err != nil {
		t.Error("got error while executing NewProxyListFromFile", err)
	}
	p := pl.FetchRandomProxy()

	u, _ := url.Parse("http://" + "1231:8000")

	if p.String() != u.String() {
		t.Error("proxies should match", p, u)
	}
}
