package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fabienzucchet/lambda-momento-extension/utils"
)

type DynamoDBClient struct {
	*dynamodb.Client
}

type Key struct {
	PK string
	SK string
}

func (k Key) String() string {
	return k.PK + k.SK

}

func InitDynamoDBClient() *DynamoDBClient {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		println(utils.PrintPrefix, "error:", err)
		return nil
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBClient{
		Client: client,
	}
}

func (c *DynamoDBClient) ReadItem(ctx context.Context, key Key, tableName string) (map[string]types.AttributeValue, error) {
	input := &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: key.PK},
			"SK": &types.AttributeValueMemberS{Value: key.SK},
		},
	}

	req, err := c.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	return req.Item, nil
}
