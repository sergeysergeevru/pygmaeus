package pygmaeus

import (
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var defaultFile string
var defaultType FileType

type FileType string

const (
	FILE_YML  FileType = "file_yml"
	FILE_JSON FileType = "file_json"
)

func init() {
	defaultFile = "config.yml"
}

func Bind(v interface{})  {
	dataByte, err := ioutil.ReadFile(defaultFile)
	if err != nil {
		logrus.Panic("can't read config")
	}
	err = yaml.Unmarshal(dataByte, v)
	if err != nil {
		panic("can't unmarshal")
	}
}

func SetDefaultFileType(t FileType)  {
	defaultType = t
}
