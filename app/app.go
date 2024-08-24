package app

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Run(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017/"))
	if err != nil {
		return err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println(err)
		}
	}()

	shortUrlDAO, err := NewUrlDAO(ctx, client)
	if err != nil {
		return err
	}

	service := NewService(shortUrlDAO)
	httpHandler := newHandler(service)

	return http.ListenAndServe(":8080", initEndpoint(httpHandler))
}

func initEndpoint(h *Handler) *echo.Echo {
	router := echo.New()

	router.POST("/shorten", WrapEndpoint(h.Shorten))
	router.GET("/:shortUrl", WrapEndpoint(h.GetFullURL))
	router.DELETE("/delete/:shortUrl", WrapEndpoint(h.Delete))
	router.POST("/update/:shortUrl", WrapEndpoint(h.Update))
	router.GET("/ping", WrapEndpoint(h.Ping))

	return router
}