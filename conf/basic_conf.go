// basic.go - basic config.

package conf

import "fmt"

type BasicConf struct {
	UrlListFile string // URLs
}

// Check checks basic config at the semantic level.
func (b *BasicConf) Check() error {
	if b.UrlListFile == "" {
		return fmt.Errorf("Empty UrlListFile")
	}

	return nil
}
