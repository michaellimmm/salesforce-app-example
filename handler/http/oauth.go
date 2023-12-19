package http

import (
	"fmt"
	"github/michaellimmm/salesforce-app-example/pkg/salesforce"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type OauthResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

type GetLoginUrlRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
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

func (h *handler) getLoginUrl(c *fiber.Ctx) error {
	req := new(GetLoginUrlRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": err})
	}

	if err := req.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": err})
	}

	res, err := h.salesforce.GetLoginUrl(c.Context(), salesforce.GetLoginUrlRequest{
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
	})
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": err})
	}
	h.logger.Info("redirect url", zap.Any("redirect_url", res))

	return c.JSON(OauthResponse{RedirectUrl: res.Url})
}

func (h *handler) oauthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": "param 'code' can not be empty"})
	}

	err := h.salesforce.ValidateAuthCode(c.Context(), code)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": "'code' does not match"})
	}

	return c.Redirect("/oauth/success")
}

func (h *handler) oauthSuccess(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}
