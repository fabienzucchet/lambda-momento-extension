package cache

import (
	"context"
	"os"
	"time"

	"github.com/fabienzucchet/lambda-momento-extension/dynamodb"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

type CacheClient struct {
	momento.CacheClient
}

var cacheName = os.Getenv("MOMENTO_CACHE_NAME")

// Initialize the Momento cache
func InitMomentoCache(token string) (*CacheClient, error) {
	credentialProvider, err := auth.FromString(token)
	if err != nil {
		return nil, err
	}

	client, err := momento.NewCacheClient(config.InRegionLatest(), credentialProvider, 600*time.Second)
	if err != nil {
		return nil, err
	}

	return &CacheClient{client}, nil
}

// Read from the cache from Momento
func (c *CacheClient) ReadCache(ctx context.Context, key dynamodb.Key) (string, error) {
	res, err := c.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String(key.String()),
	})
	if err != nil {
		return "", err
	}

	switch r := res.(type) {
	case *responses.GetHit:
		return r.ValueString(), nil
	default:
		return "", nil
	}
}

// Write to the cache
func (c *CacheClient) WriteCache(ctx context.Context, key dynamodb.Key, value string) error {
	_, err := c.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String(key.String()),
		Value:     momento.String(value),
		Ttl:       time.Duration(60 * time.Second),
	})

	return err
}
