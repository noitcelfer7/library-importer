package config

type Config struct {
	Grpc struct {
		Client struct {
			Host string `json:"host"`
			Port string `json:"port"`
		} `json:"client"`
	} `json:"grpc"`
	Http struct {
		Server struct {
			Host string `json:"host"`
			Port string `json:"port"`
		} `json:"server"`
	} `json:"http"`
}
