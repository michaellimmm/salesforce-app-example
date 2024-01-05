package http

import (
	"fmt"
	"github/michaellimmm/salesforce-app-example/pkg/salesforce"
	"net/http"

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
		return c.Render("onboard/index.onboard", fiber.Map{
			"callbackUrl": h.salesforce.GetCallbackUrl(),
		})
	})

	h.app.Post("/authorize/", func(c *fiber.Ctx) error {
		req := new(GetLoginUrlRequest)
		if err := c.BodyParser(req); err != nil {
			h.logger.Error("failed to parse body", zap.Error(err))
			return c.Render("onboard/failed.onboard", fiber.Map{"errorMessage": err.Error()})
		}

		if err := req.Validate(); err != nil {
			h.logger.Error("request body is invalid", zap.Error(err))
			return c.Render("onboard/failed.onboard", fiber.Map{"errorMessage": err.Error()})
		}

		res, err := h.salesforce.GetLoginUrl(c.Context(), salesforce.GetLoginUrlRequest{
			ClientID:     req.ClientID,
			ClientSecret: req.ClientSecret,
		})
		if err != nil {
			h.logger.Error("failed to get login url", zap.Error(err))
			return c.Render("onboard/failed.onboard", fiber.Map{"errorMessage": err.Error()})
		}

		return c.Redirect(res.Url, http.StatusTemporaryRedirect)
	})

	h.app.Get("/oauth/callback", func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code == "" {
			return c.Render("onboard/failed.onboard", fiber.Map{"errorMessage": "param 'code' can not be empty"})
		}

		err := h.salesforce.ValidateAuthCode(c.Context(), code)
		if err != nil {
			return c.Render("onboard/failed.onboard", fiber.Map{"errorMessage": "'code' does not match"})
		}

		return c.Render("onboard/success.onboard", fiber.Map{})
	})

	return h.app.Listen(addr)
}

type OauthResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

type GetLoginUrlRequest struct {
	ClientID     string `json:"client_id" form:"clientId"`
	ClientSecret string `json:"client_secret" form:"clientSecret"`
}

func (o *GetLoginUrlRequest) Validate() error {
	if o.ClientID == "" {
		return fmt.Errorf("'client_id' cannot be empty")
	}

	if o.ClientSecret == "" {
		return fmt.Errorf("'client_secret' cannot be empty")
	}

	return nil
}
