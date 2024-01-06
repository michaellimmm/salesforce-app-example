package main

import (
	"context"
	"github/michaellimmm/salesforce-app-example/db"
	"github/michaellimmm/salesforce-app-example/handlers/http"
	"github/michaellimmm/salesforce-app-example/pkg/pubsubclient"
	"github/michaellimmm/salesforce-app-example/pkg/restclient"
	"github/michaellimmm/salesforce-app-example/pkg/salesforce"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
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

	engine := html.New("./view", ".html")

	httpSrv := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "application",
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return gojson.MarshalWithOption(v, gojson.DisableHTMLEscape())
		},
	})
	httpSrv.Static("/public", "./assets")
	httpSrv.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
		Fields: []string{"url", "queryParams", "reqHeaders", "body"},
	}))
	httpSrv.Use(favicon.New(favicon.Config{
		File: "./assets/favicon.ico",
	}))

	db.Datastore = db.NewDB(context.Background(), db.WithURI("mongodb://localhost:27017"))
	db.Datastore.SelectDB("salesforce_app_db")

	restyClient := resty.New()
	restClient := restclient.NewRestClient(logger, restyClient)
	pubsubclient := pubsubclient.NewPubSubClient(logger)
	salesforceService := salesforce.NewSalesForce(logger, restClient, pubsubclient)

	logger.Info("service is running ...")

	// salesforceService.SubscribeAllLinkedToken(context.Background())

	handler := http.NewHandler(httpSrv, logger, salesforceService)
	err = handler.Serve(httpSrvPort)
	if err != nil {
		logger.Error("failed run handler", zap.Error(err))
		return
	}
}
