package config

type Config struct {
	AuthConfig AuthConfig
	FoodConfig FoodConfig
}

type AuthConfig struct {
	Key            string `env:"AUTH_KEY"`
	AuthServiceURL string `env:"AUTH_SERVICE_URL"`
}

type FoodConfig struct {
	FoodServiceAddr string `env:"FOOD_SERVICE_ADDR"`
}
