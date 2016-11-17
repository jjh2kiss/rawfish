package config

type Config struct {
	Addr           string
	Port           string
	Root           string
	ReadTimeout    int //seconds
	WriteTimeout   int //seconds
	Https          bool
	Force200Ok     bool
	Force200OkSize int
	Pemfile        string
	Process        int //count for processes
	WindowSize     int
	Rate           int
}

func (self *Config) FullAddress() string {
	return self.Addr + ":" + self.Port
}
