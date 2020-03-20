// fetcher_conf.go - Config for fetcher.

package conf

import "fmt"

type FetcherConf struct {
	CrawlTimeout int // crawl timeout, in seconds
}

// Check checks fetcher's config at the semantic level.
func (f *FetcherConf) Check() error {
	if f.CrawlTimeout <= 0 {
		return fmt.Errorf("CrawlTimeout should > 0")
	}

	return nil
}
