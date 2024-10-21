package review

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/mummumgoodboy/gateway/internal/api"
	"github.com/mummumgoodboy/gateway/internal/config"
	"github.com/mummumgoodboy/gateway/proto"
	"github.com/mummumgoodboy/verify"
)

type ReviewHandler struct {
	cfg           *config.Config
	reviewService proto.ReviewClient
	foodService   proto.RestaurantFoodClient
	verify        *verify.JWTVerifier
}

func NewReviewHandler(cfg *config.Config, reviewService proto.ReviewClient, foodService proto.RestaurantFoodClient, verifier *verify.JWTVerifier) *ReviewHandler {
	return &ReviewHandler{cfg: cfg, reviewService: reviewService, foodService: foodService, verify: verifier}
}

// CreateReview handles the creation of a review for a restaurant.
func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token", "error", err)
		return api.Unauthorized(c)
	}

	review := new(proto.ReviewRequest)
	if err := c.BodyParser(review); err != nil {
		slog.Warn("Failed to parse body", "error", err)
		return api.BadRequest(c)
	}
	review.UserId = int32(claim.UserId)

	createdReview, err := h.reviewService.CreateReview(c.Context(), review)
	if err != nil {
		slog.Warn("Failed to create review", "error", err)
		return api.ReturnError(c, err)
	}
	return c.Status(201).JSON(createdReview)
}

// GetReviewsByRestaurantId retrieves all reviews for a specific restaurant.
func (h *ReviewHandler) GetReviewsByRestaurantId(c *fiber.Ctx) error {
	_, err := h.foodService.GetRestaurantByRestaurantId(c.Context(), &proto.RestaurantIdRequest{
		Id: c.Params("restaurantId"),
	})
	if err != nil {
		slog.Warn("Failed to get restaurant", "error", err)
		return api.ReturnError(c, err)
	}

	response, err := h.reviewService.GetReviewsByRestaurantId(c.Context(), &proto.GetReviewsRequest{
		RestaurantId: c.Params("restaurantId"),
	})
	if err != nil {
		slog.Warn("Failed to retrieve reviews", "error", err)
		return api.ReturnError(c, err)
	}

	return c.JSON(response.Reviews)
}

// GetReview retrieves a specific review by its ID.
func (h *ReviewHandler) GetReview(c *fiber.Ctx) error {
	response, err := h.reviewService.GetReview(c.Context(), &proto.GetReviewRequest{
		ReviewId: c.Params("reviewId"),
	})
	if err != nil {
		slog.Warn("Failed to retrieve review", "error", err)
		return api.ReturnError(c, err)
	}

	return c.JSON(response)
}

// UpdateReview updates an existing review.
func (h *ReviewHandler) UpdateReview(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token", "error", err)
		return api.Unauthorized(c)
	}

	review := new(proto.UpdateReviewRequest)
	if err := c.BodyParser(review); err != nil {
		slog.Warn("Failed to parse body", "error", err)
		return api.BadRequest(c)
	}

	review.ReviewId = c.Params("reviewId")
	review.UserId = int32(claim.UserId)
	review.IsAdmin = claim.IsAdmin
	response, err := h.reviewService.UpdateReview(c.Context(), review)

	if err != nil {
		slog.Warn("Failed to update review", "error", err)
		return api.ReturnError(c, err)
	}

	return c.JSON(response)
}

// DeleteReview deletes a review by its ID.
func (h *ReviewHandler) DeleteReview(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token", "error", err)
		return api.Unauthorized(c)
	}

	_, err = h.reviewService.DeleteReview(c.Context(), &proto.DeleteReviewRequest{
		ReviewId: c.Params("reviewId"),
		UserId:   int32(claim.UserId),
		IsAdmin:  claim.IsAdmin,
	})
	if err != nil {
		slog.Warn("Failed to delete review", "error", err)
		return api.ReturnError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// AddFavoriteFood adds a food item to the user's list of favorites.
func (h *ReviewHandler) AddFavoriteFood(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token", "error", err)
		return api.Unauthorized(c)
	}
	foodId := c.Params("foodId")
	food, err := h.foodService.GetFoodByFoodId(c.Context(), &proto.FoodIdRequest{
		Id: foodId,
	})
	if err != nil {
		return api.ReturnError(c, err)
	}
	_, err = h.reviewService.AddFavoriteFood(c.Context(), &proto.AddFavoriteFoodRequest{
		UserId:       int32(claim.UserId),
		FoodId:       foodId,
		RestaurantId: food.RestaurantId,
	})
	if err != nil {
		slog.Warn("Failed to add favorite food", "error", err)
		return api.ReturnError(c, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

// RemoveFavoriteFood removes a food item from the user's list of favorites.
func (h *ReviewHandler) RemoveFavoriteFood(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token", "error", err)
		return api.Unauthorized(c)
	}

	foodId := c.Params("foodId")
	food, err := h.foodService.GetFoodByFoodId(c.Context(), &proto.FoodIdRequest{
		Id: foodId,
	})
	if err != nil {
		return api.ReturnError(c, err)
	}
	_, err = h.reviewService.RemoveFavoriteFood(c.Context(), &proto.RemoveFavoriteFoodRequest{
		UserId:       int32(claim.UserId),
		FoodId:       foodId,
		RestaurantId: food.RestaurantId,
	})
	if err != nil {
		slog.Warn("Failed to remove favorite food", "error", err)
		return api.ReturnError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetFavoriteFoodsByUserId retrieves a list of the user's favorite foods.
func (h *ReviewHandler) GetFavoriteFoodsByUserId(c *fiber.Ctx) error {
	claim, err := h.verify.Verify(api.GetAuthToken(c))
	if err != nil {
		slog.Warn("Failed to verify token", "error", err)
		return api.Unauthorized(c)
	}

	response, err := h.reviewService.GetFavoriteFoodsByUserId(c.Context(), &proto.GetFavoriteFoodsByUserIDRequest{
		UserId: int32(claim.UserId),
	})
	if err != nil {
		slog.Warn("Failed to retrieve favorite foods", "error", err)
		return api.ReturnError(c, err)
	}
	foodIds := []string{}
	for _, food := range response.FavoriteFoods {
		foodIds = append(foodIds, food.FoodId)
	}
	foods, err := h.foodService.GetFoodsByFoodIds(c.Context(), &proto.FoodIdsRequest{
		Ids: foodIds,
	})
	if err != nil {
		return api.ReturnError(c, err)
	}
	return c.JSON(foods)
}
