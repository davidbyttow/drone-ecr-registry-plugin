package plugin_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/davidbyttow/drone-ecr-registry-plugin/plugin"
	"github.com/drone/drone-go/plugin/registry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type fakeAccessor struct {
	api.Auth
	err error
}

func (a *fakeAccessor) GetCredentials() (*api.Auth, error) {
	if a.err != nil {
		return nil, a.err
	}
	return &a.Auth, nil
}

func TestList(t *testing.T) {
	fooAuth := api.Auth{Username: "foo", Password: "bar", ProxyEndpoint: "http://foo.com"}
	barAuth := api.Auth{Username: "bar", Password: "baz", ProxyEndpoint: "http://bar.com"}

	accessors := []plugin.RegistryAccessor{
		&fakeAccessor{Auth: fooAuth},
		&fakeAccessor{err: errors.New("failed")},
		&fakeAccessor{Auth: barAuth},
	}
	p := plugin.New(accessors, logrus.StandardLogger())
	list, err := p.List(context.Background(), &registry.Request{})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, fooAuth, api.Auth{Username: list[0].Username, Password: list[0].Password, ProxyEndpoint: list[0].Address})
	assert.Equal(t, barAuth, api.Auth{Username: list[1].Username, Password: list[1].Password, ProxyEndpoint: list[1].Address})
}
