// config_test.go - UT for config.go.

package conf

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAndCheck_NormalCase(t *testing.T) {
	confPath := "./testdata/spider.conf"

	expectConf := Config{
		Basic: BasicConf{
			"../data/url.data",
		},
		Crawler: CrawlerConf{
			1,
			1,
			8,
		},
		Fetcher: FetcherConf{
			1,
		},
		Outputer: OutputerConf{
			"../output",
			".*.(htm|html)$",
		},
	}

	conf, err := LoadAndCheck(confPath)

	assert.NoError(t, err)
	assert.Equal(t, expectConf, conf)
}

func TestLoadAndCheck_EmptyURLListFile(t *testing.T) {
	confPath := "./testdata/spider1.conf"
	_, err := LoadAndCheck(confPath)
	assert.True(t, strings.Contains(err.Error(), "Empty UrlListFile"))
}

func TestLoadAndCheck_EmptyOutputDirectory(t *testing.T) {
	confPath := "./testdata/spider2.conf"
	_, err := LoadAndCheck(confPath)
	assert.True(t, strings.Contains(err.Error(), "Empty OutputDirectory"))
}

func TestLoadAndCheck_InvalidMaxDepth(t *testing.T) {
	confPath := "./testdata/spider3.conf"
	_, err := LoadAndCheck(confPath)
	assert.True(t, strings.Contains(err.Error(), "MaxDepth should >= 0"))
}

func TestLoadAndCheck_InvalidCrawlInterval(t *testing.T) {
	confPath := "./testdata/spider4.conf"
	_, err := LoadAndCheck(confPath)
	assert.True(t, strings.Contains(err.Error(), "CrawlInterval should > 0"))
}

func TestLoadAndCheck_InvalidCrawlTimeout(t *testing.T) {
	confPath := "./testdata/spider5.conf"
	_, err := LoadAndCheck(confPath)
	assert.True(t, strings.Contains(err.Error(), "CrawlTimeout should > 0"))
}

func TestLoadAndCheck_EmptyTargetURL(t *testing.T) {
	confPath := "./testdata/spider6.conf"
	_, err := LoadAndCheck(confPath)
	assert.True(t, strings.Contains(err.Error(), "Empty TargetURL"))
}

func TestLoadAndCheck_InvalidThreadCount(t *testing.T) {
	confPath := "./testdata/spider7.conf"
	_, err := LoadAndCheck(confPath)
	assert.True(t, strings.Contains(err.Error(), "ThreadCount should > 0"))
}
