// outputer_conf.go - Config for outputer.

package conf

import "fmt"

type OutputerConf struct {
	OutputDirectory string // path of files which save result
	TargetURL       string // pattern for target URLs
}

// Check checks outputer's config at the semantic level.
func (o *OutputerConf) Check() error {
	if o.OutputDirectory == "" {
		return fmt.Errorf("Empty OutputDirectory")
	}

	if o.TargetURL == "" {
		return fmt.Errorf("Empty TargetURL")
	}

	return nil
}
