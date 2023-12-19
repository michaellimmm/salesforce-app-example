package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"

	"github/michaellimmm/salesforce-app-example/types"
)

const (
	sfAuthorizePath = "/services/oauth2/authorize"
	sfTokenPath     = "/services/oauth2/token"
	sfUserInfoPath  = "/services/oauth2/userinfo"
)

type GrantType string

const (
	GrantTypeAuthCode     GrantType = "authorization_code"
	GrantTypeRefreshToken GrantType = "refresh_token"
)

type (
	OAuth interface {
		GetToken(context.Context, TokenRequest) (TokenResponse, error)
	}

	oauth struct {
		logger *zap.Logger
	}
)

func NewOauth(logger *zap.Logger) OAuth {
	return &oauth{
		logger: logger,
	}
}

type (
	TokenRequest struct {
		GrantType    GrantType
		Code         string
		RefreshToken string
		RedirectUri  string
		CodeVerifier string
		ClientID     string
		ClientSecret string
	}

	TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Signature    string `json:"signature"`
		Scope        string `json:"scope"`
		IDToken      string `json:"id_token"`
		InstanceUrl  string `json:"instance_url"`
		ID           string `json:"id"`
		TokenType    string `json:"token_type"`
		IssuedAt     string `json:"issued_at"`
	}
)

func (t *TokenRequest) ToQueryParam() string {
	q := make(url.Values)
	q.Add("grant_type", string(t.GrantType))
	q.Add("client_id", t.ClientID)
	q.Add("client_secret", t.ClientSecret)
	q.Add("format", "json")

	switch t.GrantType {
	case GrantTypeAuthCode:
		q.Add("redirect_uri", t.RedirectUri)
		q.Add("code", t.Code)
		q.Add("code_verifier", strings.ReplaceAll(t.CodeVerifier, "%3D", "="))
	case GrantTypeRefreshToken:
		q.Add("refresh_token", t.RefreshToken)
	}

	return q.Encode()
}

func (t *TokenResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, t)
}

func (o *oauth) GetToken(ctx context.Context, param TokenRequest) (TokenResponse, error) {
	u, err := url.Parse(types.SfLoginUri + sfTokenPath)
	if err != nil {
		o.logger.Error("failed to parse url", zap.Error(err))
		return TokenResponse{}, err
	}

	u.RawQuery = param.ToQueryParam()

	o.logger.Info("url", zap.String("url", u.String()))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		o.logger.Error("failed to create request", zap.Error(err))
		return TokenResponse{}, err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		o.logger.Error("failed to get response", zap.Error(err))
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		o.logger.Error("failed to read response body", zap.Error(err))
		return TokenResponse{}, err
	}

	o.logger.Info("response", zap.Any("body", string(body)), zap.String("response code", resp.Status))

	if resp.StatusCode/100 != 2 {
		return TokenResponse{}, errors.New("failed to get")
	}

	result := TokenResponse{}
	result.Unmarshal(body)

	return result, nil
}
