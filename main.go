package main

import (
	"ShortURL/app"
	"context"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}
