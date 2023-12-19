package salesforce

import (
	"context"
	"fmt"
	"github/michaellimmm/salesforce-app-example/model"
	"github/michaellimmm/salesforce-app-example/pkg/oauth"
	"github/michaellimmm/salesforce-app-example/pkg/pubsubclient"
	"github/michaellimmm/salesforce-app-example/types"
	"github/michaellimmm/salesforce-app-example/util/crypto"
	"net/url"
	"os"

	"go.uber.org/zap"
)

const (
	sfAuthorizePath = "/services/oauth2/authorize"
	redirectPath    = "/oauth/callback"
)

type (
	Salesforce interface {
		GetLoginUrl(context.Context, GetLoginUrlRequest) (GetLoginUrlResponse, error)
		ValidateAuthCode(context.Context, string) error
	}

	salesforce struct {
		logger       *zap.Logger
		serverDomain string
		oauth        oauth.OAuth
		pubsubclient *pubsubclient.PubSubClient
	}
)

func NewSalesForce(logger *zap.Logger, oauth oauth.OAuth, pubsubclient *pubsubclient.PubSubClient) Salesforce {
	return &salesforce{
		logger:       logger,
		serverDomain: os.Getenv("HTTP_SERVER_DOMAIN"),
		oauth:        oauth,
		pubsubclient: pubsubclient,
	}
}

type (
	GetLoginUrlRequest struct {
		ClientID     string
		ClientSecret string
	}

	GetLoginUrlResponse struct {
		Url string
	}
)

func (s *salesforce) GetLoginUrl(ctx context.Context, req GetLoginUrlRequest) (GetLoginUrlResponse, error) {
	token := model.Token{
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		TokenStatus:  string(model.TokenStatusPending),
	}
	if err := token.Save(ctx); err != nil {
		return GetLoginUrlResponse{}, err
	}

	codeVerifier := token.ID.Hex()
	s.logger.Info("code_verifier", zap.String("code_verifier", codeVerifier))

	redirectUrl := s.serverDomain + redirectPath
	url, err := s.genLoginUrl(req.ClientID, redirectUrl, codeVerifier)
	if err != nil {
		return GetLoginUrlResponse{}, err
	}

	return GetLoginUrlResponse{Url: url}, nil
}

// o.serverDomain+redirectPath
func (s *salesforce) genLoginUrl(clientID, redirectUri, codeVerifier string) (string, error) {
	u, err := url.Parse(types.SfLoginUri + sfAuthorizePath)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("response_type", "code")
	q.Add("client_id", clientID)
	q.Add("redirect_uri", redirectUri)
	q.Add("code_challenge", crypto.SHA256URLEncode(codeVerifier))

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (s *salesforce) ValidateAuthCode(ctx context.Context, code string) error {
	token := model.Token{}
	tokens, err := token.FindAllByStatus(ctx, model.TokenStatusPending)
	if err != nil {
		s.logger.Error("failed to find all token by status", zap.Error(err))
		return nil
	}

	for _, t := range tokens {
		res, err := s.oauth.GetToken(ctx, oauth.TokenRequest{
			GrantType:    oauth.GrantTypeAuthCode,
			Code:         code,
			ClientID:     t.ClientID,
			ClientSecret: t.ClientSecret,
			CodeVerifier: t.ID.Hex(),
			RedirectUri:  s.serverDomain + redirectPath,
		})
		if err != nil {
			s.logger.Warn("failed to get token", zap.Error(err))
			continue
		}

		token.ClientID = t.ClientID
		token.ClientSecret = t.ClientSecret
		token.AccessToken = res.AccessToken
		token.RefreshToken = res.RefreshToken
		token.TokenStatus = string(model.TokenStatusLinked)
		token.InstanceUrl = res.InstanceUrl

		if err = token.Save(ctx); err != nil {
			s.logger.Error("failed to save token", zap.Error(err))
			return err
		}

		return nil
	}

	return fmt.Errorf("failed to get token")
}
