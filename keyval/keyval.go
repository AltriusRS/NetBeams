package keyval

import (
	"context"
	"fmt"
	"os"

	"github.com/valkey-io/valkey-go"
)

type KeyvalClient struct {
	client valkey.Client
}

func NewKeyvalClient() (*KeyvalClient, error) {
	init_address, present := os.LookupEnv("VALKEY_URI")
	if !present || init_address == "" {
		return nil, fmt.Errorf("failed to get VALKEY_URI environment variable")
	}
	valkey_password, present := os.LookupEnv("VALKEY_PASSWORD")
	if !present || valkey_password == "" {
		return nil, fmt.Errorf("failed to get VALKEY_PASSWORD environment variable")
	}

	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{init_address}, Password: valkey_password, SelectDB: 0})

	if err != nil {
		return nil, err
	}

	return &KeyvalClient{
		client: client,
	}, nil
}

func (c *KeyvalClient) Close() {
	c.client.Close()
}

func (c *KeyvalClient) Xadd(key string, id string, fields [][]string) valkey.ValkeyResult {
	builder := c.client.B().Xadd().Key(key).Id(id).FieldValue()

	for _, entities := range fields {
		if len(entities) != 2 {
			panic("Invalid field")
		}
		builder.FieldValue(entities[0], entities[1])
	}

	return c.client.Do(
		context.Background(),
		builder.Build(),
	)
}

func (c *KeyvalClient) Set(key string, value string) valkey.ValkeyResult {
	return c.client.Do(
		context.Background(),
		c.client.B().Set().Key(key).Value(value).Build(),
	)
}

func (c *KeyvalClient) Get(key string) valkey.ValkeyResult {
	return c.client.Do(
		context.Background(),
		c.client.B().Get().Key(key).Build(),
	)
}
