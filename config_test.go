package pygmaeus

import (
	"testing"
	"os"
	"reflect"
	"fmt"
)

var realLevelConfig = levelConfig{
	Root: subLevelConfig{
		SubRootOne: sub2LevelConfig{
			SubSubRoot: configTypes{
				Int64:   9223372036854775807,
				Int32:   2147483647,
				Float64: 65777.03,
				Float32: 54.05,
				String:  "lol",
				Bool:    true,
			},
		},
		SubRootTwo: configTypes{
			Int64:   -9223372036854775808,
			Int32:   -2147483648,
			Float64: -65777.03,
			Float32: -54.05,
			String:  "",
			Bool:    false,
		},
	},
}

func TestYml(t *testing.T) {
	EnableDebug(true)
	var debug levelConfig
	ReadFromFile(&debug)
	if !reflect.DeepEqual(debug, realLevelConfig) {
		t.FailNow()
	}
}

func TestJson(t *testing.T) {
	EnableDebug(true)
	SetFileType(JsonExtension)
	var debug levelConfig
	ReadFromFile(&debug)
	if !reflect.DeepEqual(debug, realLevelConfig) {
		t.FailNow()
	}
}


func TestArgs(t *testing.T) {
	var debug levelConfig
	EnableDebug(false)
	argOffset = 2
	os.Args = append(os.Args, getArgs(realLevelConfigPairs)... )
	GetFromArgs(&debug)
	t.Log(os.Args)
	if !reflect.DeepEqual(debug, realLevelConfig) {
		t.Log(debug)
		t.Log(realLevelConfig)
		t.Fail()
	}
}

func TestEnv(t *testing.T) {
	var debug levelConfig
	argOffset = 2
	setEnv(realLevelConfigPairs)
	defer disEnv(realLevelConfigPairs)
	GetFromArgs(&debug)
	if !reflect.DeepEqual(debug, realLevelConfig) {
		t.Log(debug)
		t.Log(realLevelConfig)
		t.Fail()
	}
}

var realLevelConfigPairs = [][]string{
	{"Root.SubRootOne.SubSubRoot.Int64", "9223372036854775807"},
	{"Root.SubRootOne.SubSubRoot.Int32", "2147483647"},
	{"Root.SubRootOne.SubSubRoot.Float32", "54.05"},
	{"Root.SubRootOne.SubSubRoot.Float64", "65777.03"},
	{"Root.SubRootOne.SubSubRoot.Bool", "true"},
	{"Root.SubRootOne.SubSubRoot.String", "lol"},
	{"Root.SubRootTwo.Int64", "-9223372036854775808"},
	{"Root.SubRootTwo.Int32", "-2147483648"},
	{"Root.SubRootTwo.Float32", "-54.05"},
	{"Root.SubRootTwo.Float64", "-65777.03"},
	{"Root.SubRootTwo.Bool", "false"},
}

func setEnv(pairs [][]string) {
	for _, v := range pairs {
		os.Setenv(v[0], v[1])
	}
}

func disEnv(pairs [][]string)  {
	for _, v := range pairs {
		os.Setenv(v[0], "")
	}
}

func getArgs(pairs [][]string) []string {
	a := make([]string, 0, 2*len(pairs))
	for _, v := range pairs {
		a = append(a, fmt.Sprintf("-%s", v[0]), v[1])
	}
	return a
}
