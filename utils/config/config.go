package config

type Config struct {
	JWTSecret string `yaml:"jwt_secret"`
	OllamaEndpoint string `yaml:"ollama_endpoint"`
	Model string `yaml:"model"`
	DB        DBConfig `yaml:"db"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}
