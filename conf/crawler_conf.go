// crawler_conf.go - Config for crawler.

package conf

import (
	"fmt"
)

type CrawlerConf struct {
	MaxDepth      int // max depth when crawl, depth eqauls to zero for seeds
	CrawlInterval int // crawl interval, in seconds
	ThreadCount   int // count of thread for spider
}

// Check checks crawler's config at the semantic level.
func (c *CrawlerConf) Check() error {
	if c.MaxDepth < 0 {
		return fmt.Errorf("MaxDepth should >= 0")
	}

	if c.CrawlInterval <= 0 {
		return fmt.Errorf("CrawlInterval should > 0")
	}

	if c.ThreadCount <= 0 {
		return fmt.Errorf("ThreadCount should > 0")
	}

	return nil
}
