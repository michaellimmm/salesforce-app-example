package http

import (
	"fmt"
	"github/michaellimmm/salesforce-app-example/pkg/salesforce"

	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler interface {
	Serve(string) error
}

type handler struct {
	app          *fiber.App
	logger       *zap.Logger
	salesforce   salesforce.Salesforce
	sessionStore *session.Store
}

func NewHandler(
	httpServer *fiber.App,
	logger *zap.Logger,
	salesforce salesforce.Salesforce) Handler {
	// initiate session
	sessionStore := session.New()

	return &handler{
		app:          httpServer,
		logger:       logger,
		salesforce:   salesforce,
		sessionStore: sessionStore,
	}
}

func (h *handler) Serve(addr string) error {
	h.app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("home/index", fiber.Map{})
	})

	h.app.Get("/linkage/", func(c *fiber.Ctx) error {
		c.Response().Header.Add("HX-Redirect", "/linkage/")
		return c.Render("linkage/index", fiber.Map{
			"callbackUrl": h.salesforce.GetCallbackUrl(),
		})
	})

	h.app.Post("/authorize/", func(c *fiber.Ctx) error {
		req := new(GetLoginUrlRequest)
		if err := c.BodyParser(req); err != nil {
			h.logger.Error("failed to parse body", zap.Error(err))
			return c.Render("linkage/failed", fiber.Map{"errorMessage": err.Error()})
		}

		if err := req.Validate(); err != nil {
			h.logger.Error("request body is invalid", zap.Error(err))
			return c.Render("linkage/failed", fiber.Map{"errorMessage": err.Error()})
		}

		res, err := h.salesforce.GetLoginUrl(c.Context(), salesforce.GetLoginUrlRequest{
			ClientID:     req.ClientID,
			ClientSecret: req.ClientSecret,
		})
		if err != nil {
			h.logger.Error("failed to get login url", zap.Error(err))
			return c.Render("linkage/failed", fiber.Map{"errorMessage": err.Error()})
		}

		// TODO: it seems it redirect to wrong page, fix this
		sess, err := h.sessionStore.Get(c)
		if err != nil {
			h.logger.Error("failed to get session", zap.Error(err))
			return c.Render("linkage/failed", fiber.Map{"errorMessage": err.Error()})
		}

		sess.Set("clientID", req.ClientID)
		if err := sess.Save(); err != nil {
			h.logger.Error("failed to save session", zap.Error(err))
			return c.Render("linkage/failed", fiber.Map{"errorMessage": err.Error()})

		}

		return c.Redirect(res.Url)
	})

	h.app.Get("/linkage/callback", func(c *fiber.Ctx) error {
		// TODO: it seems it redirect to wrong page, fix this
		sess, err := h.sessionStore.Get(c)
		if err != nil {
			h.logger.Error("failed to get session", zap.Error(err))
			return c.Render("linkage/failed", fiber.Map{"errorMessage": err.Error()})
		}

		code := c.Query("code")
		if code == "" {
			sess.Delete("clientID")
			return c.Render("linkage/failed", fiber.Map{"errorMessage": "param 'code' can not be empty"})
		}

		if err := h.salesforce.ValidateAuthCode(c.Context(), code); err != nil {
			sess.Delete("clientID")
			return c.Render("linkage/failed", fiber.Map{"errorMessage": "'code' does not match"})
		}

		return c.Render("linkage/success", fiber.Map{})
	})

	h.app.Get("/cdc/", func(c *fiber.Ctx) error {
		// sess, err := h.sessionStore.Get(c)
		// if err != nil {
		// 	h.logger.Error("failed to get session", zap.Error(err))
		// 	return c.Redirect("/linkage/")
		// }

		// clientID := sess.Get("clientID")
		// if clientID == "" || clientID == nil {
		// 	return c.Redirect("/linkage/")
		// }

		c.Response().Header.Add("HX-Redirect", "/cdc/")
		return c.Render("registercdc/index", fiber.Map{
			"object": salesforce.StandardObjectList,
		})
	})

	h.app.Post("/cdc/", func(c *fiber.Ctx) error {
		request := new(CDCRequest)
		_ = c.BodyParser(request)
		h.logger.Info("request", zap.Any("request", request))
		return c.Render("registercdc/success", fiber.Map{})
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

type CDCRequest struct {
	StandardObjects []string `json:"standard_objects" form:"standardObjects"`
}
