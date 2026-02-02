package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

func NewLightsailAPI(accessKeyId, secretAccessKey string) *LightsailAPI {
	client := lightsail.New(lightsail.Options{
		Region: "eu-central-1",
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKeyId,
				SecretAccessKey: secretAccessKey,
			}, nil
		}),
	})
	return &LightsailAPI{
		client: client,
	}
}

type LightsailAPI struct {
	client *lightsail.Client
}

func (api *LightsailAPI) GetServers(ctx context.Context) ([]types.Instance, error) {
	output, err := api.client.GetInstances(ctx, &lightsail.GetInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("get instances request failed: %w", err)
	}

	return output.Instances, nil
}

func (api *LightsailAPI) RebootServer(ctx context.Context, instance types.Instance) (*lightsail.RebootInstanceOutput, error) {
	return api.client.RebootInstance(ctx, &lightsail.RebootInstanceInput{
		InstanceName: instance.Name,
	})
}
