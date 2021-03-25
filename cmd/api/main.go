package main

import (
	"context"
)

func main() {
	ctx := context.Background()

	// Connect to database and load initial data
	databaseInit(&ctx)
}
