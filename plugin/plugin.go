package plugin

import (
	"context"

	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/logger"
	"github.com/drone/drone-go/plugin/registry"
	"github.com/pkg/errors"
)

// New creates a new ECR registry plugin
func New(registry string, client api.Client, logger logger.Logger) registry.Plugin {
	return &plugin{
		client:   client,
		registry: registry,
		logger:   logger,
	}
}

type plugin struct {
	client   api.Client
	registry string
	logger   logger.Logger
}

func (p *plugin) List(ctx context.Context, req *registry.Request) ([]*drone.Registry, error) {
	var list []*drone.Registry
	auth, err := p.client.GetCredentials(p.registry)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get credentials")
	}
	list = append(list, &drone.Registry{
		Address:  auth.ProxyEndpoint,
		Username: auth.Username,
		Password: auth.Password,
	})
	return list, nil
}
