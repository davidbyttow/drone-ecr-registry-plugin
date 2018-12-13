package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/davidbyttow/drone-ecr-registry-plugin/plugin"
	"github.com/drone/drone-go/plugin/registry"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type spec struct {
	Debug        bool   `envconfig:"PLUGIN_DEBUG"`
	Address      string `envconfig:"PLUGIN_ADDRESS" default:":3000"`
	Secret       string `envconfig:"PLUGIN_SECRET"`
	RegistryList string `envconfig:"ECR_REGISTRY_LIST"`
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

	logger := logrus.StandardLogger()

	var accessors []plugin.RegistryAccessor

	urls := strings.Split(spec.RegistryList, ",")
	for _, url := range urls {
		if strings.TrimSpace(url) == "" {
			continue
		}
		reg, err := api.ExtractRegistry(url)
		if err != nil {
			log.Fatal(err)
		}

		accessor := plugin.NewRegistryAccessor(reg)
		logger.Debugf("Getting credentials for %s", urls)
		if _, err = accessor.GetCredentials(); err != nil {
			logger.Errorf("unable to get credentials: %v", err)
			continue
		}

		accessors = append(accessors, accessor)
	}

	handler := registry.Handler(
		spec.Secret,
		plugin.New(
			accessors,
			logger,
		),
		logger,
	)

	logrus.Infof("Server listening on address %s", spec.Address)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Address, nil))
}
