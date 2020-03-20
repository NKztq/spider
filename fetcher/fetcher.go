// fetcher.go - fetcher can do fetch in web.

package fetcher

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/NKztq/spider/conf"
)

type Fetcher struct {
	client http.Client // client reused for fetching
}

func NewFetcher(cfg conf.FetcherConf) *Fetcher {
	client := http.Client{
		Timeout: time.Duration(cfg.CrawlTimeout) * time.Second,
	}

	return &Fetcher{client}
}

// Fetch body from URL.
func (f *Fetcher) Fetch(url string) ([]byte, error) {
	// do fetch
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest(): %v", err)
	}

	req.Header.Add("User-Agent", fakeUA())

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("url: %s, client.Get(): %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("url: %s, status code: %v", url, resp.StatusCode)
	}

	// read from URL
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}

	return body, nil
}

// Give a fake User-Agent.
func fakeUA() string {
	fakeUAs := []string{
		// chrome
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
		// IE
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)",
		// FireFox
		"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:41.0) Gecko/20100101 Firefox/41.0",
		// Safari
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_2_5 like Mac OS X) AppleWebKit/604.5.6 (KHTML, like Gecko) Version/11.0 Mobile/15D60 Safari/604.1",
	}
	rand.Seed(time.Now().Unix())
	return fakeUAs[rand.Intn(len(fakeUAs))]
}
