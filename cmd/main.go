package main

import (
	"log"
	"net/http"

	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/davidbyttow/drone-ecr-repository-plugin/plugin"
	"github.com/drone/drone-go/plugin/registry"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type spec struct {
	Debug    bool   `envconfig:"PLUGIN_DEBUG"`
	Address  string `envconfig:"PLUGIN_ADDRESS" default:":3000"`
	Secret   string `envconfig:"PLUGIN_SECRET"`
	Registry string `envconfig:"ECR_REGISTRY"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Address == "" {
		spec.Address = ":3000"
	}

	factory := &api.DefaultClientFactory{}

	reg, err := api.ExtractRegistry(spec.Registry)
	if err != nil {
		log.Fatal(err)
	}

	client := factory.NewClientFromRegion(reg.Region)

	logger := logrus.StandardLogger()

	handler := registry.Handler(
		spec.Secret,
		plugin.New(
			spec.Registry,
			client,
			logger,
		),
		logger,
	)

	logrus.Infof("Server listening on address %s", spec.Address)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Address, nil))
}
