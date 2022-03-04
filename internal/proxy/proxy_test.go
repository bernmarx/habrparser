package proxy

import (
	"net/url"
	"testing"
)

func TestNewProxyClient(t *testing.T) {
	u, _ := url.Parse("http://" + "1231:8000")

	pc := NewProxyClient(u)

	if pc.HttpClient == nil {
		t.Error("proxy client was not made")
	}
}
