// fetcher_test.go - UT for fetcher.go.

package fetcher

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/NKztq/spider/conf"
)

func TestFetch(t *testing.T) {
	conf := conf.FetcherConf{1}

	fetcher := NewFetcher(conf)

	html, err := ioutil.ReadFile("./testdata/mock.html")
	assert.NoError(t, err)

	// mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, string(html))
	}))

	// Fetch
	res, err := fetcher.Fetch(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, html, res)
}
