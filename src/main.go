package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fabienzucchet/lambda-momento-extension/cache"
	"github.com/fabienzucchet/lambda-momento-extension/dynamodb"
	"github.com/fabienzucchet/lambda-momento-extension/extension"
	"github.com/fabienzucchet/lambda-momento-extension/ipc"
	"github.com/fabienzucchet/lambda-momento-extension/utils"
)

// AWS_LAMBDA_RUNTIME_API is an environment variable that is set by the Lambda service.
// It will be accessible to the code once it runs in the Lambda execution environment.
var (
	extensionClient = extension.NewClient(os.Getenv("AWS_LAMBDA_RUNTIME_API"))
)

// The main function of the extension initializing the Momento and DynamoDB clients, the IPC and processing lambda events.
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sigs
		cancel()
		println(utils.PrintPrefix, "Received", s)
		println(utils.PrintPrefix, "Exiting")
	}()

	res, err := extensionClient.Register(ctx, utils.ExtensionName)
	if err != nil {
		panic(err)
	}
	println(utils.PrintPrefix, "Register response:", utils.PrettyPrint(res))

	// Initialize both the Momento cache client and the DynamoDB client.
	cacheClient, err := cache.InitMomentoCache(os.Getenv("MOMENTO_TOKEN"))
	if err != nil {
		panic(err)
	}
	dynamodbClient := dynamodb.InitDynamoDBClient()

	// Start the HTTP server use for the IPC.
	ipc.Start("4000", cacheClient, dynamodbClient)

	// Will block until shutdown event is received or cancelled via the context.
	processEvents(ctx)
}

// Method to process events
func processEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			println(utils.PrintPrefix, "Waiting for event...")
			res, err := extensionClient.NextEvent(ctx)
			if err != nil {
				println(utils.PrintPrefix, "Error:", err)
				println(utils.PrintPrefix, "Exiting")
				return
			}

			// Exit if we receive a SHUTDOWN event
			if res.EventType == extension.Shutdown {
				println(utils.PrintPrefix, "Received SHUTDOWN event")
				println(utils.PrintPrefix, "Exiting")
				return
			}
		}
	}
}
