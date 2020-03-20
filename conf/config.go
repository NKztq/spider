// config.go - All configs.

package conf

import (
	"fmt"

	gcfg "gopkg.in/gcfg.v1"
)

type Config struct {
	Basic    BasicConf
	Crawler  CrawlerConf
	Fetcher  FetcherConf
	Outputer OutputerConf
}

func (c *Config) load(confPath string) error {
	return gcfg.ReadFileInto(c, confPath)
}

func (c *Config) check() error {
	var err error

	err = c.Basic.Check()
	if err != nil {
		return fmt.Errorf("Basic check faild: %v", err)
	}

	err = c.Crawler.Check()
	if err != nil {
		return fmt.Errorf("Crawler check faild: %v", err)
	}

	err = c.Fetcher.Check()
	if err != nil {
		return fmt.Errorf("Fetcher check faild: %v", err)
	}

	err = c.Outputer.Check()
	if err != nil {
		return fmt.Errorf("Outputer check faild: %v", err)
	}

	return nil
}

// LoadAndCheck loads config from file.
//
// Param:
//	- confPath: file path of config.
//
// Returns:
//	- (SpiderConf, err msg).
func LoadAndCheck(confPath string) (Config, error) {
	var conf Config
	var err error

	err = conf.load(confPath)
	if err != nil {
		return conf, fmt.Errorf("load(): %v", err)
	}

	err = conf.check()
	if err != nil {
		return conf, fmt.Errorf("check(): %v", err)
	}

	return conf, nil
}
