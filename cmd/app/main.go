package main

import (
	"context"
	"fmt"
	"github.com/pablodev/s3-test/internal/entity"
)

func main() {
	ctx := context.Background()
	app, err := InitializeApp(ctx)
	if err != nil {
		panic(err)
	}

	// Example: Upload
	file := entity.YamlFile{ID: "test-1", Version: 1, Expression: "1+2"}
	err = app.Upload(ctx, file)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("Application initialized and running.")
}
