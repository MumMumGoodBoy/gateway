package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/mummumgoodboy/gateway/internal/config"
	"github.com/mummumgoodboy/gateway/internal/handler/auth"
	"github.com/mummumgoodboy/gateway/internal/handler/food"
	"github.com/mummumgoodboy/gateway/internal/handler/recommend"
	"github.com/mummumgoodboy/gateway/internal/route"
	"github.com/mummumgoodboy/gateway/proto"
	"github.com/mummumgoodboy/verify"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg config.Config
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	verifier, err := verify.NewJWTVerifier(cfg.AuthConfig.Key)
	if err != nil {
		log.Fatal(err)
	}

	foodServiceConn, err := grpc.NewClient(cfg.FoodConfig.FoodServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	foodService := proto.NewRestaurantFoodClient(foodServiceConn)

	recommendServiceConn, err := grpc.NewClient(cfg.RecommendConfig.RecommendServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	recommendService := proto.NewRecommendServiceClient(recommendServiceConn)

	authHandler := auth.NewAuthHandler(&cfg)
	foodHandler := food.NewFoodHandler(&cfg, foodService, verifier)
	recommendHandler := recommend.NewRecommendHandler(&cfg, foodService, recommendService, verifier)

	router := route.Route{
		AuthHandler:      authHandler,
		FoodHandler:      foodHandler,
		RecommendHandler: recommendHandler,
	}

	app := fiber.New()

	router.Apply(app)

	log.Println("Gateway is running on port 3000")
	app.Listen(":3000")
}
