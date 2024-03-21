package cmd

import (
	"context"
	"log"
	app "testTaskHezzl/internal/app"
)

func main() {
	ctx := context.Background()
	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err = a.Run(); err != nil {
		log.Fatal(err)
	}
}
