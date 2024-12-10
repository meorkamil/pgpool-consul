package model

type Config struct {
	Consul struct {
		Addr     string `yaml:"addr"`
		Services struct {
			Name     string `yaml:"name"`
			Addr     string `yaml:"addr"`
			Port     int    `yaml:"port"`
			Interval string `yaml:"interval"`
			Timeout  string `yaml:"timeout"`
		} `yaml:"services"`
	} `yaml:"consul"`
	Pgpool struct {
		Listen      string `yaml:"listen"`
		Pcppassfile string `yaml:"pcppassfile"`
		Pcpport     int    `yaml:"pcpport"`
		Pcpuser     string `yaml:"pcpuser"`
		Id          string `yaml:"id"`
	} `yaml:"pgpool"`
	Global struct {
		Interval int `yaml:"interval"`
	} `yaml:"global"`
	Version string
}
