package ipc

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/fabienzucchet/lambda-momento-extension/cache"
	"github.com/fabienzucchet/lambda-momento-extension/dynamodb"
	"github.com/fabienzucchet/lambda-momento-extension/utils"
)

// Start begins running the sidecar
func Start(port string, cacheClient *cache.CacheClient, dynamodbClient *dynamodb.DynamoDBClient) {
	go startHTTPServer(port, cacheClient, dynamodbClient)
}

// IPC server handling requests. Fetches the requested item from the cache or DynamoDB and returns it.
// If caching is disabled, it will always fetch from DynamoDB.
func startHTTPServer(port string, cacheClient *cache.CacheClient, dynamodbClient *dynamodb.DynamoDBClient) {
	http.HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) {
		pk := r.URL.Query().Get("pk")
		sk := r.URL.Query().Get("sk")
		key := dynamodb.Key{PK: pk, SK: sk}
		println(utils.PrintPrefix, "Received request for pk/sk:", pk, sk)
		if os.Getenv("CACHING_DISABLED") == "true" {
			println(utils.PrintPrefix, "Caching is disabled, reading from DynamoDB...")
			value, err := fetchItemFromDynamoDB(r.Context(), key, dynamodbClient)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Write(value)
		} else {
			value, err := fetchItemIfCached(r.Context(), key, cacheClient, dynamodbClient)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Write(value)
		}
	})
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
	println(utils.PrintPrefix, "Server started on port", port)
}

// Function implementing the caching logic
func fetchItemIfCached(context context.Context, key dynamodb.Key, cacheClient *cache.CacheClient, dynamodbClient *dynamodb.DynamoDBClient) ([]byte, error) {
	value, err := cacheClient.ReadCache(context, key)
	if err != nil {
		return nil, err
	}

	if value == "" {
		println("Cache miss for key", key.String(), "Reading from DynamoDB...")

		item, err := dynamodbClient.ReadItem(context, key, os.Getenv("DYNAMODB_TABLE_NAME"))

		if err != nil || item == nil {
			return nil, err
		}

		err = cacheClient.WriteCache(context, key, utils.PrettyPrint(item))
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(&item)
		if err != nil {
			return nil, err
		}

		return b, nil
	}

	println("Cache hit for key", key.String(), "Returning cached value...")
	return []byte(value), nil
}

// Function fetching straight from DynamoDB
func fetchItemFromDynamoDB(context context.Context, key dynamodb.Key, dynamodbClient *dynamodb.DynamoDBClient) ([]byte, error) {
	item, err := dynamodbClient.ReadItem(context, key, os.Getenv("DYNAMODB_TABLE_NAME"))
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(&item)
	if err != nil {
		return nil, err
	}

	return b, nil
}
