package search

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mummumgoodboy/gateway/internal/api"
	"github.com/mummumgoodboy/gateway/internal/config"
)

// Act as a gateway to the search service
type SearchHandler struct {
	cfg *config.Config
}

func NewSearchHandler(cfg *config.Config) *SearchHandler {
	return &SearchHandler{cfg: cfg}
}

func (h *SearchHandler) SearchRestaurants(c *fiber.Ctx) error {
	searchQuery := c.Query("search")
	offset := c.Query("offset")
	limit := c.Query("limit")

	url := fmt.Sprintf("%s/search/restaurants?search=%s&offset=%s&limit=%s", h.cfg.SearchConfig.SearchServiceAddr, searchQuery, offset, limit)
	return api.HandleRedirect(url, c)
}

func (h *SearchHandler) SearchFoods(c *fiber.Ctx) error {
	searchQuery := c.Query("search")
	offset := c.Query("offset")
	limit := c.Query("limit")
	maxPrice := c.Query("maxPrice")
	minPrice := c.Query("minPrice")

	url := fmt.Sprintf("%s/search/foods?search=%s&offset=%s&limit=%s&maxPrice=%s&minPrice=%s", h.cfg.SearchConfig.SearchServiceAddr, searchQuery, offset, limit, maxPrice, minPrice)
	return api.HandleRedirect(url, c)
}
