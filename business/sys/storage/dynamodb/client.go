package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type client struct {
	Project string
	Stage   string
	Region  string
	Client  *dynamodb.Client
}

// NewClient returns client with dynamodb client integrated with config.
func NewClient(project, stage, region, profile string) (*client, error) {
	defaultLocalCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}

	// Using the Config value, create the DynamoDB client
	return &client{
		Project: project,
		Stage:   stage,
		Region:  region,
		Client:  dynamodb.NewFromConfig(defaultLocalCfg),
	}, nil
}

// tableName returns default name of the table by using base entrypoint.
func (c *client) tableName(base string) *string {
	return aws.String(fmt.Sprintf("%s-%s-%s", base, c.Project, c.Stage))
}

