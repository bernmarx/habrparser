package proxyfetcher

import (
	"io/ioutil"
	"math/rand"
	"net/url"
	"strings"
)

type ProxyList struct {
	URL []*url.URL
}

func NewProxyListFromFile(filePathAbs string) (ProxyList, error) {
	var proxies ProxyList

	plb, err := ioutil.ReadFile(filePathAbs)

	if err != nil {
		return proxies, err
	}

	pls := string(plb)

	pls = strings.Trim(pls, "\n")

	ps := strings.Split(pls, ";")
	for _, p := range ps {
		u, err := url.Parse("http://" + p)
		if err != nil {
			return proxies, err
		}

		proxies.URL = append(proxies.URL, u)
	}

	return proxies, nil
}

func (p *ProxyList) FetchRandomProxy() *url.URL {
	n := rand.Intn(len(p.URL))
	return p.URL[n]
}
