package config

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}

type configLoader interface {
	Load(filePath string)
}

func LoadConfig(config configLoader, filePath string) {
	config.Load(filePath)
}
