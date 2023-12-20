package restclient

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

const (
	userInfoEndpoint = "/services/oauth2/userinfo"
)

type UserInfo interface {
	GetUserInfo(ctx context.Context, instanceUrl, accessToken string) (UserInfoResponse, error)
}

type UserInfoResponse struct {
	OrgID string `json:"organization_id"`
}

func (u *UserInfoResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, u)
}

func (r *restClient) GetUserInfo(ctx context.Context, instanceUrl, accessToken string) (UserInfoResponse, error) {
	result := UserInfoResponse{}
	url := instanceUrl + userInfoEndpoint
	resp, err := r.client.R().SetResult(&result).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		Get(url)
	if err != nil {
		r.logger.Error("failed to get user info", zap.Error(err))
		return UserInfoResponse{}, err
	}

	r.logger.Info("response",
		zap.Any("body", string(resp.Body())),
		zap.String("response code", resp.Status()))

	return result, nil
}
