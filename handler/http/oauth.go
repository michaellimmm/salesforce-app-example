package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type OauthResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

func (h *handler) generateLoginUrl(c *fiber.Ctx) error {
	redirectUrl, err := h.oauthService.GenerateLoginUrl()
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": "failed get redirect url"})
	}
	h.logger.Info("redirect url", zap.Any("redirect_url", redirectUrl))

	return c.JSON(OauthResponse{RedirectUrl: redirectUrl})
}
