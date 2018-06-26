Simple configuration library for binding configuration to structure from config.yml (had to be at working directory), environment variables and cli arguments.
A priority of value for structure is next: config.yml (the lowest priority), env vars, cli args (the major priority).

Using example:
~~~
var Config struct{
    Server struct {
        Host string `yaml:"host"`
        Port string `yaml:"port"`
    }  `yaml:"server"`
}

pygmaeus.Bind(&config)
~~~

Main aim of this library take experience with reflect Go library. Also I express my point of view of configuration library.
