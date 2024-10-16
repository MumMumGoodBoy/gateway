package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mummumgoodboy/gateway/internal/api"
	"github.com/mummumgoodboy/gateway/internal/config"
)

// Act as a gateway to the auth service
type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	return api.HandleRedirect(h.cfg.AuthConfig.AuthServiceURL+"/auth/login", c)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	return api.HandleRedirect(h.cfg.AuthConfig.AuthServiceURL+"/auth/register", c)
}
