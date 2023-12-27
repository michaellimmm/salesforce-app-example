package http

import (
	"github/michaellimmm/salesforce-app-example/pkg/salesforce"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler interface {
	Serve(string) error
}

type handler struct {
	app        *fiber.App
	logger     *zap.Logger
	salesforce salesforce.Salesforce
}

func NewHandler(
	httpServer *fiber.App,
	logger *zap.Logger,
	salesforce salesforce.Salesforce) Handler {
	return &handler{
		app:        httpServer,
		logger:     logger,
		salesforce: salesforce,
	}
}

func (h *handler) Serve(addr string) error {
	h.app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("onboard/index.onboard", fiber.Map{})
	})

	h.app.Get("/callback", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"callback_url": h.salesforce.GetCallbackUrl(),
		})
	})
	h.app.Post("/oauth", h.getLoginUrl)
	h.app.Get("/oauth/callback", h.oauthCallback)
	h.app.Get("/oauth/success", h.oauthSuccess)

	return h.app.Listen(addr)
}
