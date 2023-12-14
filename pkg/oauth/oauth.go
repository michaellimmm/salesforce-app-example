package oauth

import (
	"context"
	"encoding/json"
	"github/michaellimmm/salesforce-app-example/util/crypto"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"go.uber.org/zap"
)

const (
	sfOAuthUri      = "https://login.salesforce.com"
	sfAuthorizePath = "/services/oauth2/authorize"
	sfTokenPath     = "/services/oauth2/token"
	sfUserInfoPath  = "/services/oauth2/userinfo"
	redirectPath    = "/oauth/callback"
)

type GrantType string

const (
	GrantTypeAuthCode     GrantType = "authorization_code"
	GrantTypeRefreshToken GrantType = "refresh_token"
)

type (
	OAuth interface {
		GenerateLoginUrl() (string, error)
		GetToken(context.Context, TokenRequest) error
	}

	oauth struct {
		logger       *zap.Logger
		serverDomain string
		clientId     string
		clientSecret string
		clientCode   string
	}
)

func NewOauth(logger *zap.Logger) OAuth {
	return &oauth{
		logger:       logger,
		serverDomain: os.Getenv("HTTP_SERVER_DOMAIN"),
		clientId:     os.Getenv("CLIENT_ID"),
		clientSecret: os.Getenv("CLIENT_SECRET"),
		clientCode:   os.Getenv("CLIENT_CODE"),
	}
}

func (o *oauth) GenerateLoginUrl() (string, error) {
	u, err := url.Parse(sfOAuthUri + sfAuthorizePath)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("response_type", "code")
	q.Add("client_id", o.clientId)
	q.Add("redirect_uri", o.serverDomain+redirectPath)
	q.Add("code_challenge", crypto.SHA256URLEncode(o.clientCode))

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (o *oauth) GetToken(ctx context.Context, param TokenRequest) error {
	u, err := url.Parse(sfOAuthUri + sfTokenPath)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Add("grant_type", "authorization_code")
	q.Add("client_id", o.clientId)
	q.Add("client_secret", o.clientSecret)
	q.Add("format", "json")

	switch param.GrantType {
	case GrantTypeAuthCode:
		q.Add("redirect_uri", o.serverDomain+redirectPath)
		q.Add("code", param.Code)
		q.Add("code_verifier", o.clientCode)
	case GrantTypeRefreshToken:
		q.Add("refresh_token", param.RefreshToken)
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	result := TokenResponse{}
	result.Unmarshal(body)

	o.logger.Info("body", zap.Any("body", result))

	return nil
}

type TokenRequest struct {
	GrantType    GrantType
	Code         string
	RefreshToken string
}

type TokenResponse struct {
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

func (t *TokenResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, t)
}
