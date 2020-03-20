// parser.go - parse html content.

package parser

import (
	"net/url"

	"golang.org/x/net/html"
)

// Parse a html page, find deeper URLs recursively.
//
// Params:
//	- n: html node.
//	- u: base URL for relative URLs in this node.
//	- deeperURLs: output param, used for saving deeper URLs.
func Parse(n *html.Node, u *url.URL, deeperURLs *[]*url.URL) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}

			// parse URL and filter out invalid URL
			rawURL, err := url.Parse(a.Val)
			if err != nil {
				continue
			}

			*(deeperURLs) = append(*(deeperURLs), u.ResolveReference(rawURL))
		}
	}

	// tail recursion
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		Parse(c, u, deeperURLs)
	}
}
