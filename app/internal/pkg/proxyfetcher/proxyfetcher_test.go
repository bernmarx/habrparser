package proxyfetcher

import (
	"net/url"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadProxyFromFile(t *testing.T) {
	proxiesPath, err := filepath.Abs("proxies_test.txt")
	assert.Nil(t, err)

	got, err := LoadProxyFromFile(proxiesPath)
	assert.Nil(t, err)

	var expected []Proxy
	u1, _ := url.Parse("http://123:123")
	u2, _ := url.Parse("http://4214:1234")
	expected = append(expected, Proxy{u1})
	expected = append(expected, Proxy{u2})

	assert.Equal(t, expected, got)
}

func TestFetchRandomProxy(t *testing.T) {
	proxiesPath, err := filepath.Abs("proxies_test.txt")
	assert.Nil(t, err)

	proxies, err := LoadProxyFromFile(proxiesPath)
	assert.Nil(t, err)

	rp := FetchRandomProxy(proxies)

	for _, p := range proxies {
		if p == rp {
			return
		}
	}
	t.Error("random proxy does not match any proxy from proxy list")
}
