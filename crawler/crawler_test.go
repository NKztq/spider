// crawler_test.go - UT for crawler.go.

package crawler

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/NKztq/spider/conf"
	"github.com/NKztq/spider/parser"
)

// implement for Fetcher
type mockFetcher struct{}

// implement for Outputer
type mockOutputer struct {
	outputDirectory string
}

func (m *mockFetcher) Fetch(url string) ([]byte, error) {
	ret := map[string][]byte{
		"http://www.baidu.com":  []byte("test"),
		"http://www.baidu1.com": []byte("test1"),
		"http://www.baidu2.com": []byte("test2"),
	}

	return ret[url], nil
}

func (m *mockOutputer) OutputFile(fileName string, content []byte) error {
	if _, err := os.Stat(m.outputDirectory); os.IsNotExist(err) {
		os.Mkdir(m.outputDirectory, os.ModePerm)
	}

	f, _ := os.Create(path.Join(m.outputDirectory, fileName))
	defer f.Close()

	f.Write(content)

	return nil
}

func TestRunOnce(t *testing.T) {
	guard := monkey.Patch(parser.Parse, func(n *html.Node, u *url.URL, d *[]*url.URL) {
		u1, _ := url.Parse("http://www.baidu1.com")
		u2, _ := url.Parse("http://www.baidu2.com")
		*d = append(*d, u1)
		*d = append(*d, u2)
	})
	defer guard.Unpatch()

	outputDirectory := "./testoutput"

	// param
	cfg := conf.CrawlerConf{
		MaxDepth:      1,
		CrawlInterval: 1,
		ThreadCount:   8,
	}
	seeds := []string{"http://www.baidu.com"}
	fetcher := &mockFetcher{}
	outputer := &mockOutputer{outputDirectory}

	// new
	crawler := NewCrawler(cfg, seeds, fetcher, outputer)

	// run
	crawler.RunOnce()

	// should crawl www.baidu.com, www.baidu1.com, www.baidu2.com
	data, err := ioutil.ReadFile("./testoutput/http%3A%2F%2Fwww.baidu.com")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), data)
	data1, err := ioutil.ReadFile("./testoutput/http%3A%2F%2Fwww.baidu1.com")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test1"), data1)
	data2, err := ioutil.ReadFile("./testoutput/http%3A%2F%2Fwww.baidu2.com")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test2"), data2)

	// delete testoutput
	assert.NoError(t, os.RemoveAll(outputDirectory))
}

func TestRunOnce_DepthZero(t *testing.T) {
	guard := monkey.Patch(parser.Parse, func(n *html.Node, u *url.URL, d *[]*url.URL) {
		u1, _ := url.Parse("http://www.baidu1.com")
		u2, _ := url.Parse("http://www.baidu2.com")
		*d = append(*d, u1)
		*d = append(*d, u2)
	})
	defer guard.Unpatch()

	outputDirectory := "./testoutput1"

	// param
	cfg := conf.CrawlerConf{
		MaxDepth:      0,
		CrawlInterval: 1,
		ThreadCount:   8,
	}
	seeds := []string{"http://www.baidu.com"}
	fetcher := &mockFetcher{}
	outputer := &mockOutputer{outputDirectory}

	// new
	crawler := NewCrawler(cfg, seeds, fetcher, outputer)

	// run
	crawler.RunOnce()

	// should only crawl www.baidu.com
	data, err := ioutil.ReadFile("./testoutput1/http%3A%2F%2Fwww.baidu.com")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), data)
	_, err = ioutil.ReadFile("./testoutput1/http%3A%2F%2Fwww.baidu1.com")
	assert.True(t, strings.Contains(err.Error(), "no such file or directory"))
	_, err = ioutil.ReadFile("./testoutput1/http%3A%2F%2Fwww.baidu2.com")
	assert.True(t, strings.Contains(err.Error(), "no such file or directory"))

	// delete testoutput
	assert.NoError(t, os.RemoveAll(outputDirectory))
}

func TestLimitFrequency(t *testing.T) {
	guard := monkey.Patch(parser.Parse, func(n *html.Node, u *url.URL, d *[]*url.URL) {
		u1, _ := url.Parse("http://www.baidu1.com")
		u2, _ := url.Parse("http://www.baidu2.com")
		*d = append(*d, u1)
		*d = append(*d, u2)
	})
	defer guard.Unpatch()

	outputDirectory := "./testoutput2"

	// param
	cfg := conf.CrawlerConf{
		MaxDepth:      1,
		CrawlInterval: 2,
		ThreadCount:   8,
	}
	seeds := []string{"http://www.baidu1.com/test.html"}
	fetcher := &mockFetcher{}
	outputer := &mockOutputer{outputDirectory}

	// new
	crawler := NewCrawler(cfg, seeds, fetcher, outputer)

	now := time.Now().Unix()
	crawler.RunOnce()
	after := time.Now().Unix()

	if after-now < 2 {
		t.Errorf("limitFrequency failed")
	}

	// delete testoutput
	assert.NoError(t, os.RemoveAll(outputDirectory))
}

func TestRunOnce_ShortTaskQueue(t *testing.T) {
	mem := taskQueueLength
	taskQueueLength = 1
	defer func() {
		taskQueueLength = mem
	}()

	guard := monkey.Patch(parser.Parse, func(n *html.Node, u *url.URL, d *[]*url.URL) {
		u1, _ := url.Parse("http://www.baidu1.com")
		u2, _ := url.Parse("http://www.baidu2.com")
		*d = append(*d, u1)
		*d = append(*d, u2)
	})
	defer guard.Unpatch()

	outputDirectory := "./testoutput3"

	// param
	cfg := conf.CrawlerConf{
		MaxDepth:      1,
		CrawlInterval: 1,
		ThreadCount:   8,
	}
	seeds := []string{"http://www.baidu.com", "http://www.baidu1.com"}
	fetcher := &mockFetcher{}
	outputer := &mockOutputer{outputDirectory}

	// new
	crawler := NewCrawler(cfg, seeds, fetcher, outputer)

	// run
	crawler.RunOnce()

	// should crawl www.baidu.com, www.baidu1.com, www.baidu2.com
	data, err := ioutil.ReadFile("./testoutput3/http%3A%2F%2Fwww.baidu.com")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), data)
	data1, err := ioutil.ReadFile("./testoutput3/http%3A%2F%2Fwww.baidu1.com")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test1"), data1)
	data2, err := ioutil.ReadFile("./testoutput3/http%3A%2F%2Fwww.baidu2.com")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test2"), data2)

	// delete testoutput
	assert.NoError(t, os.RemoveAll(outputDirectory))
}
