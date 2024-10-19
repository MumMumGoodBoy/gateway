package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mummumgoodboy/gateway/internal/handler/auth"
	"github.com/mummumgoodboy/gateway/internal/handler/food"
)

type Route struct {
	AuthHandler *auth.AuthHandler
	FoodHandler *food.FoodHandler
}

func (r *Route) Apply(f fiber.Router) {
	auth := f.Group("/auth")
	auth.Post("/login", r.AuthHandler.Login)
	auth.Post("/register", r.AuthHandler.Register)
	auth.Get("/me", r.AuthHandler.GetMe)
	auth.Put("/me", r.AuthHandler.UpdateProfile)
	auth.Patch("/me/password", r.AuthHandler.ChangePassword)

	food := f.Group("/food")
	food.Get("/:id", r.FoodHandler.GetFood)
	food.Post("/", r.FoodHandler.CreateFood)
	food.Put("/:id", r.FoodHandler.UpdateFood)
	food.Delete("/:id", r.FoodHandler.DeleteFood)

	restaurant := f.Group("/restaurant")
	restaurant.Get("/", r.FoodHandler.GetRestaurants)
	restaurant.Get("/:id", r.FoodHandler.GetRestaurant)
	restaurant.Post("/", r.FoodHandler.CreateRestaurant)
	restaurant.Put("/:id", r.FoodHandler.UpdateRestaurant)
	restaurant.Delete("/:id", r.FoodHandler.DeleteRestaurant)
	restaurant.Get("/:id/foods", r.FoodHandler.GetFoodsByRestaurantId)
}
