package restclient

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

const (
	sfLoginUri      = "https://login.salesforce.com"
	sfAuthorizePath = "/services/oauth2/authorize"
	sfTokenPath     = "/services/oauth2/token"
	sfUserInfoPath  = "/services/oauth2/userinfo"
)

type GrantType string

const (
	GrantTypeAuthCode     GrantType = "authorization_code"
	GrantTypeRefreshToken GrantType = "refresh_token"
)

type OAuth interface {
	GetToken(context.Context, TokenRequest) (TokenResponse, error)
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

func (r *restClient) GetToken(ctx context.Context, param TokenRequest) (TokenResponse, error) {
	u, err := url.Parse(sfLoginUri + sfTokenPath)
	if err != nil {
		r.logger.Error("failed to parse url", zap.Error(err))
		return TokenResponse{}, err
	}

	u.RawQuery = param.ToQueryParam()

	r.logger.Info("url", zap.String("url", u.String()))

	result := TokenResponse{}
	resp, err := r.client.R().SetResult(&result).Post(u.String())
	if err != nil {
		r.logger.Error("failed to get token", zap.Error(err))
		return TokenResponse{}, err
	}

	r.logger.Info("response",
		zap.Any("body", string(resp.Body())),
		zap.String("response code", resp.Status()))

	return result, nil
}
