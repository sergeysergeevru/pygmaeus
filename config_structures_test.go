package pygmaeus

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
	SubRootOne sub2LevelConfig `yaml:"level_2_1"`
	SubRootTwo configTypes     `yaml:"level_2_2"`
}

type sub2LevelConfig struct {
	SubSubRoot configTypes `yaml:"level_3_1"`
}

type levelConfig struct {
	Root subLevelConfig `yaml:"level_1_1"`
}

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

var realConfigLocal = levelConfig{
	Root: subLevelConfig{
		SubRootOne: sub2LevelConfig{
			SubSubRoot: configTypes{
				Int64:   986655,
				Int32:   2147,
				Float64: 5434.55,
				Float32: 657.03,
				String:  "local config",
				Bool:    false ,
			},
		},
		SubRootTwo: configTypes{
			Int64:   -922337203,
			Int32:   -214,
			Float64: -5324.15,
			Float32: -65777.03,
			String:  "",
			Bool:    true,
		},
	},
}