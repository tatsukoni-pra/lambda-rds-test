package main

import (
	"app/dataStore"
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (string, error) {
	dataStore.UpdateDBData()
	return "", nil
}

func main() {
	lambda.Start(HandleRequest)
}
