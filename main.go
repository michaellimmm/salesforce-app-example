package main

import (
	"context"
	"github/michaellimmm/salesforce-app-example/db"
	"github/michaellimmm/salesforce-app-example/handler/http"
	"github/michaellimmm/salesforce-app-example/pkg/oauth"
	"github/michaellimmm/salesforce-app-example/pkg/pubsubclient"
	"github/michaellimmm/salesforce-app-example/pkg/salesforce"
	"log"
	"os"

	gojson "github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Error loading init logger")
	}
	defer logger.Sync()

	httpSrvPort := os.Getenv("HTTP_SERVER_PORT")

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	engine := html.New("./view", ".html")

	httpSrv := fiber.New(fiber.Config{
		Views: engine,
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return gojson.MarshalWithOption(v, gojson.DisableHTMLEscape())
		},
	})
	httpSrv.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
		Fields: []string{"url", "queryParams", "reqHeaders", "body"},
	}))

	db.Datastore = db.NewDB(context.Background(), db.WithURI("mongodb://localhost:27017"))
	db.Datastore.SelectDB("salesforce_app_db")

	oauthService := oauth.NewOauth(logger)
	pubsubclient := pubsubclient.NewPubSubClient(logger)
	salesforceService := salesforce.NewSalesForce(logger, oauthService, pubsubclient)

	logger.Info("service is running ...")

	handler := http.NewHandler(httpSrv, logger, salesforceService)
	err = handler.Serve(httpSrvPort)
	if err != nil {
		logger.Error("failed run handler", zap.Error(err))
		return
	}
}
