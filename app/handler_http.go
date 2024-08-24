package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
)

type Handler struct {
	service *Service
}

type UrlRequest struct {
	Url     string `json:"url"`
	TTLDays int    `json:"ttlDays"`
}

func newHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Shorten(ctx echo.Context) (interface{}, error) {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Printf("Failed reading body request: %v", err)
		return nil, err
	}

	var reqJson UrlRequest
	if err = json.Unmarshal(body, &reqJson); err != nil {
		log.Printf("Failed unmarshaling body: %v", err)
		return nil, err
	}

	if _, err = url.ParseRequestURI(reqJson.Url); err != nil {
		log.Printf("Failed parse request URI: %v", err)
		return nil, err
	}

	return h.service.Shorten(ctx.Request().Context(), reqJson.Url, reqJson.TTLDays)
}

func (h *Handler) Ping(ctx echo.Context) (interface{}, error) {
	return nil, nil
}

func (h *Handler) Update(ctx echo.Context) (interface{}, error) {
	id := ctx.Param("shortUrl")
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, err
	}
	
	var reqJson UrlRequest
	err = json.Unmarshal(body, &reqJson)
	if err != nil {
		return nil, err
	}

	if _, err = url.ParseRequestURI(reqJson.Url); err != nil {
		return nil, err
	}

	return h.service.Update(ctx.Request().Context(), id, reqJson.Url, reqJson.TTLDays)
}

func (h *Handler) GetFullURL(ctx echo.Context) (interface{}, error) {
	shortUrl := ctx.Param("shortUrl")
	log.Println(shortUrl)
	return h.service.GetFullURL(ctx.Request().Context(), shortUrl)
}

func (h *Handler) Delete(ctx echo.Context) (interface{}, error) {
	id := ctx.Param("shortUrl")
	log.Println(id)
	return nil, h.service.Delete(ctx.Request().Context(), id)
}

type EndpointHandler func(ctx echo.Context) (interface{}, error)

func WrapEndpoint(h EndpointHandler) func (echo.Context) error {
	fn := func(ctx echo.Context, h EndpointHandler) error {
		result, err := h(ctx)
		if err != nil {
			return err
		}

		data, err := json.Marshal(result)
		if err != nil {
			return err
		}

		_, err = ctx.Response().Write(data)
		return err
	}
	return func(ctx echo.Context) error{
		err := fn(ctx, h)
		if err != nil {
			log.Println(err.Error())
			ctx.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
		return nil
	}
}