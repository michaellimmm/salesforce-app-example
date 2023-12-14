package http

import (
	"github/michaellimmm/salesforce-app-example/pkg/oauth"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler interface {
	Serve(string) error
}

type handler struct {
	app          *fiber.App
	logger       *zap.Logger
	oauthService oauth.OAuth
}

func NewHandler(httpServer *fiber.App, logger *zap.Logger, oauthService oauth.OAuth) Handler {
	return &handler{
		app:          httpServer,
		logger:       logger,
		oauthService: oauthService,
	}
}

func (h *handler) Serve(addr string) error {
	h.app.Get("/oauth", h.generateLoginUrl)

	h.app.Get("/oauth/callback", func(c *fiber.Ctx) error {

		code := c.Query("code")
		if code == "" {
			return c.Status(http.StatusUnprocessableEntity).
				JSON(fiber.Map{"error": "param 'code' can not be empty"})
		}

		err := h.oauthService.GetToken(c.Context(), oauth.TokenRequest{
			GrantType: oauth.GrantTypeAuthCode,
			Code:      code,
		})
		if err != nil {
			h.logger.Error("error", zap.Error(err))
		}

		return c.Redirect("/oauth/thanks")
	})

	h.app.Get("/oauth/thanks", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	return h.app.Listen(addr)
}
