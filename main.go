package main

import (
	"context"
	"fmt"
	"log"
	"workerpool/run"
)

func main() {
	ctx := context.Background()

	if err := run.Execute(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Run `go test -race -v ./...` to validate.")

}
