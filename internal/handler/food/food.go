package food

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/mummumgoodboy/gateway/internal/api"
	"github.com/mummumgoodboy/gateway/internal/config"
	"github.com/mummumgoodboy/gateway/proto"
	"github.com/mummumgoodboy/verify"
)

type FoodHandler struct {
	cfg *config.Config

	foodService proto.RestaurantFoodClient
	verify      *verify.JWTVerifier
}

func NewFoodHandler(cfg *config.Config, foodService proto.RestaurantFoodClient, verifier *verify.JWTVerifier) *FoodHandler {
	return &FoodHandler{cfg: cfg, foodService: foodService, verify: verifier}
}

func (h *FoodHandler) GetFood(c *fiber.Ctx) error {
	food, err := h.foodService.GetFoodByFoodId(c.Context(), &proto.FoodIdRequest{
		Id: c.Params("id"),
	})
	if err != nil {
		return api.ReturnError(c, err)
	}

	return c.JSON(food)
}

func (h *FoodHandler) CreateFood(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token",
			"error", err,
		)
		return api.Unauthorized(c)
	}

	if !claim.IsAdmin {
		slog.Warn("User is not admin",
			"user", claim.UserId,
		)
		return api.Forbidden(c)
	}

	food := new(proto.Food)
	if err := c.BodyParser(food); err != nil {
		slog.Warn("Failed to parse body",
			"error", err)
		return api.BadRequest(c)
	}

	food, err = h.foodService.CreateFood(c.Context(), food)
	if err != nil {
		slog.Warn("Failed to create food",
			"error", err)
		return api.ReturnError(c, err)
	}

	return c.JSON(food)
}

func (h *FoodHandler) UpdateFood(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token",
			"error", err,
		)
		return api.Unauthorized(c)
	}

	if !claim.IsAdmin {
		slog.Warn("User is not admin",
			"user", claim.UserId,
		)
		return api.Forbidden(c)
	}

	food := new(proto.Food)
	if err := c.BodyParser(food); err != nil {
		slog.Warn("Failed to parse body",
			"error", err)
		return api.BadRequest(c)
	}

	food.Id = c.Params("id")

	food, err = h.foodService.UpdateFood(c.Context(), food)
	if err != nil {
		slog.Warn("Failed to update food",
			"error", err)
		return api.ReturnError(c, err)
	}

	return c.JSON(food)
}

func (h *FoodHandler) DeleteFood(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token",
			"error", err,
		)
		return api.Unauthorized(c)
	}

	if !claim.IsAdmin {
		slog.Warn("User is not admin",
			"user", claim.UserId,
		)
		return api.Forbidden(c)
	}

	_, err = h.foodService.DeleteFood(c.Context(), &proto.FoodIdRequest{
		Id: c.Params("id"),
	})
	if err != nil {
		slog.Warn("Failed to delete food",
			"error", err)
		return api.ReturnError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *FoodHandler) GetFoodsByRestaurantId(c *fiber.Ctx) error {
	_, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token",
			"error", err,
		)
		return api.Unauthorized(c)
	}
	foods, err := h.foodService.GetFoodsByRestaurantId(c.Context(), &proto.RestaurantIdRequest{
		Id: c.Params("id"),
	})
	if err != nil {
		return api.ReturnError(c, err)
	}

	return c.JSON(foods)
}

func (h *FoodHandler) GetRestaurants(c *fiber.Ctx) error {
	restaurants, err := h.foodService.GetRestaurants(c.Context(), new(proto.Empty))
	if err != nil {
		return api.ReturnError(c, err)
	}

	return c.JSON(restaurants)
}

func (h *FoodHandler) GetRestaurant(c *fiber.Ctx) error {
	restaurant, err := h.foodService.GetRestaurantByRestaurantId(c.Context(), &proto.RestaurantIdRequest{
		Id: c.Params("id"),
	})
	if err != nil {
		return api.ReturnError(c, err)
	}

	return c.JSON(restaurant)
}

func (h *FoodHandler) CreateRestaurant(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token",
			"error", err,
		)
		return api.Unauthorized(c)
	}

	if !claim.IsAdmin {
		slog.Warn("User is not admin",
			"user", claim.UserId,
		)
		return api.Forbidden(c)
	}

	req := new(proto.CreateRestaurantRequest)
	if err := c.BodyParser(req); err != nil {
		slog.Warn("Failed to parse body",
			"error", err)
		return api.BadRequest(c)
	}

	restaurant, err := h.foodService.CreateRestaurant(c.Context(), req)
	if err != nil {
		slog.Warn("Failed to create restaurant",
			"error", err)
		return api.ReturnError(c, err)
	}

	return c.JSON(restaurant)
}

func (h *FoodHandler) UpdateRestaurant(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token",
			"error", err,
		)
		return api.Unauthorized(c)
	}

	if !claim.IsAdmin {
		slog.Warn("User is not admin",
			"user", claim.UserId,
		)
		return api.Forbidden(c)
	}

	req := new(proto.Restaurant)
	if err := c.BodyParser(req); err != nil {
		slog.Warn("Failed to parse body",
			"error", err)
	}

	req.Id = c.Params("id")

	restaurant, err := h.foodService.UpdateRestaurants(c.Context(), req)
	if err != nil {
		slog.Warn("Failed to update restaurant",
			"error", err)
		return api.ReturnError(c, err)
	}

	return c.JSON(restaurant)
}

func (h *FoodHandler) DeleteRestaurant(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token",
			"error", err,
		)
		return api.Unauthorized(c)
	}

	if !claim.IsAdmin {
		slog.Warn("User is not admin",
			"user", claim.UserId,
		)
		return api.Forbidden(c)
	}

	_, err = h.foodService.DeleteRestaurant(c.Context(), &proto.RestaurantIdRequest{
		Id: c.Params("id"),
	})
	if err != nil {
		slog.Warn("Failed to delete restaurant",
			"error", err)
		return api.ReturnError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
