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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {

	}

	defer func () {
		if err := client.Disconnect(ctx); err != nil {
			log.Println(err)
		}
	}()

	httpHandler := newHandler()

	return http.ListenAndServe(":8080", initEndpoint(httpHandler))
}

func initEndpoint(h *Handler) *echo.Echo {
	router := echo.New()

	router.POST("/shorten", h.Shorten)
	router.GET("/:shortURL", h.GetFullURL)
	router.DELETE("/:shortURL", h.Delete)
	router.POST("/update/:shortUrl", h.Update)
	router.GET("/ping", h.Ping)

	return router
}