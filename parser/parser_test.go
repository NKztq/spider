// parser_test.go - UT for parse.go.

package parser

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestParse(t *testing.T) {
	page, err := ioutil.ReadFile("./testdata/mock.html")
	assert.NoError(t, err)

	node, err := html.Parse(bytes.NewReader(page))
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "http://www.baidu.com/test/test/test", nil)
	assert.NoError(t, err)

	deeperURLs := []*url.URL{}
	Parse(node, req.URL, &deeperURLs)

	// deeper URLs ordered like:
	//
	//	-mock1.html
	//		-mock5.html
	//			-mock8.html
	//	-mock2.html
	//		-mock6.html
	//	-mock3.html
	//		-mock7.html
	//	-mock4.html
	//	-../path.html
	//	-../../path.html
	//	-//host/path
	//	-http://h.o.s.t/test
	//	-http://h.o.s.t/test/
	//	-http://h.o.s.t/test/#more
	expectedURLs := []string{"http://1.1.1.1:80/mock1.html", "http://1.1.1.1:80/mock5.html", "http://1.1.1.1:80/mock8.html", "http://1.1.1.1:80/mock2.html", "http://1.1.1.1:80/mock6.html", "http://1.1.1.1:80/mock3.html", "http://1.1.1.1:80/mock7.html", "http://1.1.1.1:80/mock4.html", "http://www.baidu.com/test/path.html", "http://www.baidu.com/path.html", "http://host/path", "http://h.o.s.t/test", "http://h.o.s.t/test/", "http://h.o.s.t/test/#more"}

	actualURLs := []string{}
	for _, u := range deeperURLs {
		actualURLs = append(actualURLs, u.String())
	}

	assert.Equal(t, expectedURLs, actualURLs)
}
