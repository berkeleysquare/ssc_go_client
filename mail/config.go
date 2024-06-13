package mail

type Config struct {
	Email struct {
		Server        string `yaml:"server"`
		Port          int    `yaml:"port"`
		Email         string `yaml:"from"`
		Password      string `yaml:"password"`
		Authorization string `yaml:"authorization"`
	}
	Message struct {
		Subject  string   `yaml:"subject"`
		Template string   `yaml:"template"`
		To       []string `yaml:"to"`
	}
	Params struct {
		Endpoint string `yaml:"endpoint"`
	}
}
