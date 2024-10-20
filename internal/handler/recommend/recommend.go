package recommend

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/mummumgoodboy/gateway/internal/api"
	"github.com/mummumgoodboy/gateway/internal/config"
	"github.com/mummumgoodboy/gateway/proto"
	"github.com/mummumgoodboy/verify"
)

type RecommendHandler struct {
	cfg *config.Config

	foodService      proto.RestaurantFoodClient
	recommendService proto.RecommendServiceClient
	verify           *verify.JWTVerifier
}

func NewRecommendHandler(cfg *config.Config, foodService proto.RestaurantFoodClient, recommendService proto.RecommendServiceClient, verify *verify.JWTVerifier) *RecommendHandler {
	return &RecommendHandler{
		cfg:              cfg,
		foodService:      foodService,
		recommendService: recommendService,
		verify:           verify,
	}
}

func (h *RecommendHandler) GetRecommend(c *fiber.Ctx) error {
	token := api.GetAuthToken(c)

	userID := 0

	if token != "" {
		claim, err := h.verify.Verify(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		userID = int(claim.UserId)
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	withNoDelay := c.QueryBool("no_delay", false)

	// Get recommend food
	recommendFood, err := h.recommendService.GetFoodRecommendations(c.Context(),
		&proto.GetRecommendationsRequest{
			UserId:  int64(userID),
			Limit:   int32(limit),
			Offset:  int32(offset),
			NoDelay: withNoDelay,
		})
	if err != nil {
		slog.Warn("Error while getting recommendation",
			"err", err,
		)
		return api.InternalError(c)
	}

	res, err := h.foodService.GetFoodsByFoodIds(c.Context(), &proto.FoodIdsRequest{
		Ids: recommendFood.ItemIds,
	})
	if err != nil {
		slog.Warn("Error while getting food by ids",
			"err", err,
		)
		return api.InternalError(c)
	}

	return c.JSON(res.Foods)
}
