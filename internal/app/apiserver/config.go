package apiserver

type config struct {
	BindAddr string `toml:"bind_addr"`
}

func NewConfig() *config {
	return &config{
		BindAddr: ":8080",
	}
}