package main

import (
	"context"
	"log"
	"wiki-export/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.New(ctx)
	if err != nil {
		log.Fatalf("[ERROR] while initializing app: %s", err)
	}

	if err := a.Run(); err != nil {
		log.Fatalf("[ERROR] while running app: %s", err)
	}
}
