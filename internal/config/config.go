package config

type Config struct {
	AuthConfig      AuthConfig
	FoodConfig      FoodConfig
	RecommendConfig RecommendConfig
	ReviewConfig    ReviewConfig
}

type AuthConfig struct {
	Key            string `env:"AUTH_KEY"`
	AuthServiceURL string `env:"AUTH_SERVICE_URL"`
}

type FoodConfig struct {
	FoodServiceAddr string `env:"FOOD_SERVICE_ADDR"`
}

type RecommendConfig struct {
	RecommendServiceAddr string `env:"RECOMMEND_SERVICE_ADDR"`
}

type ReviewConfig struct {
	ReviewServiceAddr string `env:"REVIEW_SERVICE_ADDR"`
}

type SearchConfig struct {
	SearchServiceAddr string `env:"SEARCH_SERVICE_ADDR"`
}
