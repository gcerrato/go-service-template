package main

import (
	"context"
	_ "embed"
	"flag"
	"net"
	"os"
	"strings"

	// Import flowchartsman/swaggerui package
	"github.com/gcerrato/go-service-template/api"
	"github.com/gcerrato/go-service-template/database/ent"
	"github.com/gcerrato/go-service-template/internal/repos"
	"github.com/gcerrato/go-service-template/internal/services"
	"github.com/gcerrato/swaggerui-echo"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	validationMiddleware "github.com/oapi-codegen/echo-middleware"
	"github.com/sirupsen/logrus"
)

// Database setup function
func databaseSetup() *ent.Client {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("Error loading .env file")
	}
	logrus.Info("host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT") + " user=" + os.Getenv("POSTGRES_USER") + " dbname=" + os.Getenv("POSTGRES_DB") + " password=" + os.Getenv("POSTGRES_PASSWORD"))

	logrus.Info("open db")
	dbClient, err := ent.Open("postgres", "sslmode=disable host="+os.Getenv("POSTGRES_HOST")+" port="+os.Getenv("POSTGRES_PORT")+" user="+os.Getenv("POSTGRES_USER")+" dbname="+os.Getenv("POSTGRES_DB")+" password="+os.Getenv("POSTGRES_PASSWORD"))
	if err != nil {
		logrus.Error("error ent open", "err", err.Error())
	}

	logrus.Info("create schema")
	if err := dbClient.Schema.Create(context.Background()); err != nil {
		logrus.Error("failed creating schema resources", "err", err)
	}
	return dbClient
}

func main() {
	port := flag.String("port", "3000", "Port for test HTTP server")
	flag.Parse()

	// Load OpenAPI V3 specification (Swagger JSON) file
	swagger, err := api.GetSwagger()
	if err != nil {
		logrus.Error("Error loading swagger spec", "err", err.Error())
		os.Exit(1)
	}

	db := databaseSetup()
	defer db.Close()

	// Clear out the servers array in the swagger spec
	swagger.Servers = nil

	todoRepo := repos.NewTodoRepo(db)
	todoService := services.NewTodoService(todoRepo)

	apiHandler := api.NewServerHandler(*todoService)
	println(apiHandler)

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.CORS())
	api.RegisterHandlers(e, apiHandler)
	validatorOptions := &validationMiddleware.Options{
		Skipper: func(c echo.Context) bool {
			// Skip validation for any path that starts with /swagger/
			return strings.HasPrefix(c.Request().URL.Path, "/swagger/")
		},
	}

	// Apply the OpenAPI validation middleware with the options
	e.Use(validationMiddleware.OapiRequestValidatorWithOptions(swagger, validatorOptions))

	e.GET("/swagger/*", swaggerui.EchoHandler("/swagger", api.Spec))

	logrus.Info("starting http server...", "port", *port)

	// Start the server
	if err := e.Start(net.JoinHostPort("0.0.0.0", *port)); err != nil {
		logrus.Error("server error", "err", err)
		os.Exit(1)
	}
}
