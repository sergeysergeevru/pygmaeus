package pygmaeus

import (
	"testing"
	"reflect"
	"os"
	"log"
)

type config struct {
	Name  string
	Year  int64
	Month int
	Server struct {
		Host   string
		Port   int
		Tokens []string
	}
	Debug bool
}

var realConfig = config{
	Name:  "test",
	Year:  922337203685477580,
	Month: 3,
	Server: struct {
		Host   string
		Port   int
		Tokens []string
	}{Host: "localhost", Port: 8888, Tokens: []string{"one", "two", "three"}},
	Debug: true,
}

type configTypes struct {
	Int64   int64
	Int32   int64
	Float32 float32
	Float64 float64
	Bool    bool
	String  string
}


type subLevelConfig struct {
	Level21 sub2LevelConfig `yaml:"level_2_1"`
	Level22 configTypes `yaml:"level_2_2"`
}

type sub2LevelConfig struct {
	Level31 configTypes `yaml:"level_3_1"`
}

type levelConfig struct {
	Level11 subLevelConfig `yaml:"level_1_1"`
}

var  realLevelConfig = levelConfig{
	Level11:subLevelConfig{
		Level21:sub2LevelConfig{
			Level31: configTypes{
				Int64: 9223372036854775807,
				Int32: 2147483647,
				Float64: 65777.03,
				Float32: 54.05,
				String: "lol",
				Bool: true,
			},
		},
		Level22: configTypes{
			Int64: -9223372036854775808,
			Int32: -2147483648,
			Float64: -65777.03,
			Float32: -54.05,
			String: "",
			Bool: false,
		},
	},
}


func TestYml(t *testing.T) {
	var debug levelConfig
	Bind(&debug)
	log.Println(debug, realLevelConfig)
	if !reflect.DeepEqual(debug, realLevelConfig) {
		t.FailNow()
	}
}

func TestArgs(t *testing.T) {
	var debug config
	os.Args = append(os.Args,
		"-Name", "1234",
		//"-Server.Host", "localhost",
		"-Server.Port", "111",
		"-Month", "03",
		"-Year", "922337203685477580")
	os.Setenv("Server.Host", "yam")
	//var st string
	GetFromArgs(&debug)
	//t.Log(os.Args)
	t.Log(debug)
	//if !reflect.DeepEqual(debug, realConfig) {
	//	t.Log(debug, realConfig)
	//	t.Fail()
	//}
}
