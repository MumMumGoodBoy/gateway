package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mummumgoodboy/gateway/internal/handler/auth"
	"github.com/mummumgoodboy/gateway/internal/handler/food"
	"github.com/mummumgoodboy/gateway/internal/handler/recommend"
	"github.com/mummumgoodboy/gateway/internal/handler/review"
)

type Route struct {
	AuthHandler      *auth.AuthHandler
	FoodHandler      *food.FoodHandler
	RecommendHandler *recommend.RecommendHandler
	ReviewHandler    *review.ReviewHandler
}

func (r *Route) Apply(f fiber.Router) {
	auth := f.Group("/auth")
	auth.Post("/login", r.AuthHandler.Login)
	auth.Post("/register", r.AuthHandler.Register)
	auth.Get("/me", r.AuthHandler.GetMe)
	auth.Put("/me", r.AuthHandler.UpdateProfile)
	auth.Patch("/me/password", r.AuthHandler.ChangePassword)

	food := f.Group("/food")
	food.Get("/:foodId", r.FoodHandler.GetFood)
	food.Put("/:foodId", r.FoodHandler.UpdateFood)
	food.Delete("/:foodId", r.FoodHandler.DeleteFood)

	restaurant := f.Group("/restaurant")
	restaurant.Get("/", r.FoodHandler.GetRestaurants)
	restaurant.Get("/:restaurantId", r.FoodHandler.GetRestaurant)
	restaurant.Post("/", r.FoodHandler.CreateRestaurant)
	restaurant.Put("/:restaurantId", r.FoodHandler.UpdateRestaurant)
	restaurant.Delete("/:restaurantId", r.FoodHandler.DeleteRestaurant)

	restaurant.Get("/:restaurantId/food", r.FoodHandler.GetFoodsByRestaurantId)

	foodRecommend := f.Group("/food-recommend")
	foodRecommend.Get("/", r.RecommendHandler.GetRecommend)
	restaurant.Post("/:restaurantId/food", r.FoodHandler.CreateFood)

	restaurant.Post("/:restaurantId/review", r.ReviewHandler.CreateReview)
	restaurant.Get("/:restaurantId/review", r.ReviewHandler.GetReviewsByRestaurantId)

	review := f.Group("/review")
	review.Get("/:reviewId", r.ReviewHandler.GetReview)
	review.Put("/:reviewId", r.ReviewHandler.UpdateReview)
	review.Delete("/:reviewId", r.ReviewHandler.DeleteReview)

	favorite := f.Group("/favorite")
	favorite.Post("/:foodId", r.ReviewHandler.AddFavoriteFood)
	favorite.Delete("/:foodId", r.ReviewHandler.RemoveFavoriteFood)
	favorite.Get("/", r.ReviewHandler.GetFavoriteFoodsByUserId)
}
