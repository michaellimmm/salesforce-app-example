package oauth_test

import (
	"github/michaellimmm/salesforce-app-example/pkg/oauth"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestGenerateLoginUrl(t *testing.T) {
	serverDomain := "http://localhost:8080"
	clientId := "secretbutnotsecret"

	os.Setenv("HTTP_SERVER_DOMAIN", serverDomain)
	os.Setenv("CLIENT_ID", clientId)
	os.Setenv("CLIENT_SECRET", "clientsecret")
	os.Setenv("CLIENT_CODE", "anothercode")

	oauth := oauth.NewOauth(zap.NewNop())
	actual, err := oauth.GenerateLoginUrl()
	if err != nil {
		t.Fatalf("Failed, actual %v want %v", err, nil)
	}

	if actual == "" {
		t.Fatal("Failed, data cannot be empty")
	}

	if !strings.Contains(actual, "https://login.salesforce.com/services/oauth2/authorize") {
		t.Fatalf("Failed, actual: %v doesn't contain correct endpoint", actual)
	}

	if !strings.Contains(actual, "client_id="+clientId) {
		t.Fatalf("Failed, actual: %v doesn't contain correct endpoint", actual)
	}

	if !strings.Contains(actual, `redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Foauth%2Fcallback`) {
		t.Fatalf("Failed, actual: %v doesn't contain correct endpoint", actual)
	}

	if !strings.Contains(actual, `response_type=code`) {
		t.Fatalf("Failed, actual: %v doesn't contain correct endpoint", actual)
	}
}
