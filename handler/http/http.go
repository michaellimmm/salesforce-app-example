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
	h.app.Post("/oauth", h.getLoginUrl)
	h.app.Get("/oauth/callback", h.oauthCallback)
	h.app.Get("/oauth/success", h.oauthSuccess)

	return h.app.Listen(addr)
}
