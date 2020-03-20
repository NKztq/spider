// outputer.go - output content to file.

package outputer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/baidu/go-lib/log"

	"github.com/NKztq/spider/conf"
)

const (
	fileNameMaxLength = 255 // in Bytes, Linux&MacOS's max file name length
	md5HashLength     = 32
)

type Outputer struct {
	OutputDirectory string
	Pattern         *regexp.Regexp
}

func NewOutputer(cfg conf.OutputerConf) (*Outputer, error) {
	pattern, err := regexp.Compile(cfg.TargetURL)
	if err != nil {
		return nil, fmt.Errorf("url: %s, regexp.Compile(): %v", cfg.TargetURL, err)
	}

	return &Outputer{cfg.OutputDirectory, pattern}, nil
}

// Output content into file whose path is joined by Outputer's outputDirectory and fileName.
// FileNames that match failed will not output.
func (o *Outputer) OutputFile(fileName string, content []byte) error {
	if !o.Pattern.MatchString(fileName) {
		log.Logger.Info("OutputFile(): url: %s match failed", fileName)
		return nil
	}

	fileName = hashLongFileName(fileName)

	var err error

	_, err = os.Stat(o.OutputDirectory)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("directory: %s, os.Stat(): %v", o.OutputDirectory, err)
	}

	// mkdir if not exist
	if os.IsNotExist(err) {
		err = os.Mkdir(o.OutputDirectory, os.ModePerm)
		if err != nil {
			return fmt.Errorf("directory: %s, os.Mkdir(): %v", o.OutputDirectory, err)
		}
	}

	fp := path.Join(o.OutputDirectory, fileName)

	f, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("os.Create(): ap: %s, err: %v", fp, err)
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf("write to file: %s failed, err: %v", fp, err)
	}

	log.Logger.Info("OutputFile(): url: %s output successfully", fileName)

	return nil
}

// For file names that longer than fileNameMaxLength,
// do md5 hash for [(fileNameMaxLength - md5HashLength):] of the file name,
// append hash result to [:(fileNameMaxLength - md5HashLength)] of the file name
// as new file name.
func hashLongFileName(fileName string) string {
	if len(fileName) <= fileNameMaxLength {
		return fileName
	}

	reserve := fileName[:(fileNameMaxLength - md5HashLength)]
	needHash := fileName[(fileNameMaxLength - md5HashLength):]

	md5Inst := md5.New()
	md5Inst.Write([]byte(needHash))
	hashed := hex.EncodeToString(md5Inst.Sum([]byte("")))

	return reserve + hashed
}
