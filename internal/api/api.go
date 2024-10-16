package api

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ErrorResp struct {
	Message string `json:"message"`
}

func InternalError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResp{
		Message: "Internal server error",
	})
}

func Unauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(ErrorResp{
		Message: "Unauthorized",
	})
}

func Forbidden(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(ErrorResp{
		Message: "Forbidden",
	})
}

func BadRequest(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(ErrorResp{
		Message: "Bad request",
	})
}

func ReturnError(c *fiber.Ctx, err error) error {
	slog.Warn("Error in handling request",
		"error", err,
	)
	return InternalError(c)
}

func ReturnResp(c *fiber.Ctx, resp *http.Response) error {
	body := new(bytes.Buffer)
	_, err := body.ReadFrom(resp.Body)
	if err != nil {
		slog.Warn("Error reading response body",
			"error", err)
		return ReturnError(c, err)
	}
	c.Set(fiber.HeaderContentType, resp.Header.Get(fiber.HeaderContentType))
	return c.
		Status(resp.StatusCode).
		SendStream(body)
}

func RedirectRequest(url string, c *fiber.Ctx) (*http.Response, error) {
	req, err := http.NewRequest(c.Method(), url, bytes.NewReader(c.Body()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Copy headers from the original request to the new request
	for k, vs := range c.GetReqHeaders() {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	// Copy query parameters from the original request to the new request
	q := req.URL.Query()
	for k, v := range c.Queries() {
		q.Add(k, v)
	}

	return http.DefaultClient.Do(req)
}

func HandleRedirect(url string, c *fiber.Ctx) error {
	resp, err := RedirectRequest(url, c)
	if err != nil {
		return ReturnError(c, err)
	}
	defer resp.Body.Close()

	return ReturnResp(c, resp)
}

func GetAuthToken(c *fiber.Ctx) string {
	token, found := strings.CutPrefix(c.Get("Authorization"), "Bearer ")
	if !found {
		return ""
	}

	return token
}
