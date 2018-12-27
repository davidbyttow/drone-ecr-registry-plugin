package plugin

import (
	"context"
	"time"
	"fmt"
	"strings"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/logger"
	"github.com/drone/drone-go/plugin/registry"
	"github.com/pkg/errors"
	vault "github.com/hashicorp/vault/api"
	"github.com/kelseyhightower/envconfig"
)

var factory = &api.DefaultClientFactory{}

type RegistryAccessor interface {
	GetCredentials() (*api.Auth, error)
}

type accessor struct {
	registry *api.Registry
	client   api.Client
	auth     *api.Auth
	expiry   time.Time
}

type gcrlist struct {
	GCRRegistryList string `envconfig:"GCR_REGISTRY_LIST"`
	GCRVaultPath string `envconfig:"GCR_VAULT_PATH"`
}

func NewRegistryAccessor(registry *api.Registry) *accessor {
	return &accessor{
		registry: registry,
	}
}

func (ra *accessor) GetCredentials() (*api.Auth, error) {
	if ra.client == nil {
		ra.client = factory.NewClientFromRegion(ra.registry.Region)
	}
	if ra.auth == nil || ra.expiry.Before(time.Now()) {
		ra.auth = nil
		var err error
		ra.auth, err = ra.client.GetCredentialsByRegistryID(ra.registry.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't get credentials for registry: %s", ra.registry.ID)
		}
		// NOTE: These expire in 12 hours, but this is OK.
		ra.expiry = time.Now().Add(time.Hour)
	}
	return ra.auth, nil
}

type plugin struct {
	accessors []RegistryAccessor
	logger    logger.Logger
}

// New creates a new ECR registry plugin
func New(accessors []RegistryAccessor, logger logger.Logger) registry.Plugin {
	return &plugin{
		accessors: accessors,
		logger:    logger,
	}
}

func (p *plugin) List(ctx context.Context, req *registry.Request) ([]*drone.Registry, error) {
	var list []*drone.Registry

	for _, acc := range p.accessors {
		auth, err := acc.GetCredentials()
		if err != nil {
			p.logger.Errorf("unable to fetch credentials for registry: %v", err)
			continue
		}
		list = append(list, &drone.Registry{
			Address:  auth.ProxyEndpoint,
			Username: auth.Username,
			Password: auth.Password,
		})
	}

	GCR := new(gcrlist)
	err := envconfig.Process("", GCR)
	if err != nil {
		fmt.Println(err)
	}
	urls := strings.Split(GCR.GCRRegistryList, ",")
	gcr_key := get_vault_key(GCR.GCRVaultPath)
	for _, host := range urls {
		list = append(list, &drone.Registry{
			Address: host,
			Username: "_json_key",
			Password: gcr_key,
		})
	}
	return list, nil
}

func get_vault_key(path string) string {
	config := &vault.Config{}
	config.ReadEnvironment()
	Client, _ := vault.NewClient(config)
	l := Client.Logical()
	secret, _ := l.Read(path)
	datamap := secret.Data
	data := datamap["value"].(string)
	return data
}
