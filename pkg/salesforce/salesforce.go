package salesforce

import (
	"context"
	"errors"
	"fmt"
	"github/michaellimmm/salesforce-app-example/model"
	"github/michaellimmm/salesforce-app-example/pkg/pubsubclient"
	"github/michaellimmm/salesforce-app-example/pkg/restclient"
	"github/michaellimmm/salesforce-app-example/util/crypto"
	"net/url"
	"os"

	"go.uber.org/zap"
)

const (
	sfLoginUri      = "https://login.salesforce.com"
	sfAuthorizePath = "/services/oauth2/authorize"

	redirectPath = "/oauth/callback"
)

type (
	Salesforce interface {
		GetLoginUrl(context.Context, GetLoginUrlRequest) (GetLoginUrlResponse, error)
		ValidateAuthCode(context.Context, string) error
	}

	salesforce struct {
		logger       *zap.Logger
		serverDomain string
		restClient   restclient.RestClient
		pubsubclient *pubsubclient.PubSubClient
	}
)

func NewSalesForce(logger *zap.Logger, restClient restclient.RestClient, pubsubclient *pubsubclient.PubSubClient) Salesforce {
	return &salesforce{
		logger:       logger,
		serverDomain: os.Getenv("HTTP_SERVER_DOMAIN"),
		restClient:   restClient,
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
	}
	if err := token.FindByClientIDAndClientSecret(ctx); err != nil {
		if !errors.Is(err, model.ErrDataNotFound) {
			return GetLoginUrlResponse{}, err
		}

		token.TokenStatus = string(model.TokenStatusPending)
		if err := token.Save(ctx); err != nil {
			return GetLoginUrlResponse{}, err
		}
	}

	codeVerifier := token.ID.Hex()
	redirectUrl := s.serverDomain + redirectPath
	url, err := s.genLoginUrl(req.ClientID, redirectUrl, codeVerifier)
	if err != nil {
		return GetLoginUrlResponse{}, err
	}

	return GetLoginUrlResponse{Url: url}, nil
}

func (s *salesforce) genLoginUrl(clientID, redirectUri, codeVerifier string) (string, error) {
	u, err := url.Parse(sfLoginUri + sfAuthorizePath)
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

	for i := 0; i < len(tokens); i++ {
		newToken := tokens[i]

		req := restclient.TokenRequest{
			GrantType:    restclient.GrantTypeAuthCode,
			Code:         code,
			ClientID:     newToken.ClientID,
			ClientSecret: newToken.ClientSecret,
			CodeVerifier: newToken.ID.Hex(),
			RedirectUri:  s.serverDomain + redirectPath,
		}
		tokenResp, err := s.restClient.GetToken(ctx, req)
		if err != nil {
			s.logger.Warn("failed to get token", zap.Error(err))
			continue
		}

		userInfoResp, err := s.restClient.GetUserInfo(ctx, tokenResp.InstanceUrl, tokenResp.AccessToken)
		if err != nil {
			s.logger.Error("failed get user info", zap.Error(err))
			return err
		}

		newToken.AccessToken = tokenResp.AccessToken
		newToken.RefreshToken = tokenResp.RefreshToken
		newToken.TokenStatus = string(model.TokenStatusLinked)
		newToken.InstanceUrl = tokenResp.InstanceUrl
		newToken.OrgID = userInfoResp.OrgID

		if err = newToken.Update(ctx); err != nil {
			s.logger.Error("failed to save token", zap.Error(err))
			return err
		}

		return nil
	}

	return fmt.Errorf("auth code is not valid")
}
