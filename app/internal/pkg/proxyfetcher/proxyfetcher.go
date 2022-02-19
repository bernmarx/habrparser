package proxyfetcher

import (
	"io/ioutil"
	"math/rand"
	"net/url"
	"strings"
)

type Proxy struct {
	URL *url.URL
}

func FetchRandomProxy(proxies []Proxy) Proxy {
	n := rand.Intn(len(proxies))
	return proxies[n]
}

func LoadProxyFromFile(absPath string) ([]Proxy, error) {
	var proxies []Proxy

	plb, err := ioutil.ReadFile(absPath)

	if err != nil {
		return nil, err
	}

	pls := string(plb)

	ps := strings.Split(pls, ";")
	for _, p := range ps {
		u, err := url.Parse("http://" + p)
		if err != nil {
			return nil, err
		}

		proxies = append(proxies, Proxy{u})
	}

	return proxies, nil
}
