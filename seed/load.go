

// load.go - Load seeds.

package seed

import (
	"encoding/json"
	"io/ioutil"
)

// Load loads seeds from file, will drop duplicate seeds.
/*
Seeds in file like:
[
     "http://www.baidu.com",
     "http://www.sina.com.cn",
     ...
   ]
*/
func Load(seedPath string) (seeds []string, err error) {
	rawData, err := ioutil.ReadFile(seedPath)
	if err != nil {
		return
	}

	rawSeeds := make([]string, 0)
	err = json.Unmarshal(rawData, &rawSeeds)
	if err != nil {
		return
	}

	// drop duplicate seeds
	seedMap := make(map[string]bool)
	for _, seed := range rawSeeds {
		seedMap[seed] = true
	}

	seeds = make([]string, 0, len(seedMap))
	for seed := range seedMap {
		seeds = append(seeds, seed)
	}

	return
}
