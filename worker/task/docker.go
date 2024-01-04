package task

import (
	"context"
	"log"

	"github.com/docker/docker/client"
)

var cli *client.Client
var ctx context.Context

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error initializing Docker client: %v", err)
		ctx = context.Background()
	}
}
