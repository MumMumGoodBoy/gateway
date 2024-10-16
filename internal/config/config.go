package config

type Config struct {
	AuthConfig AuthConfig
}

type AuthConfig struct {
	Key            string `env:"AUTH_KEY"`
	AuthServiceURL string `env:"AUTH_SERVICE_URL"`
}
