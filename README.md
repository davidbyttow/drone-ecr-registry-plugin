# drone-ecr-repository-plugin

Drone extension to support ECR repositories when pulling images

## Building
To build and push to your repository
```
$ make linux-amd64-image
$ docker tag drone-ecr-registry-plugin:<sha> your/repo:latest
$ docker push your/repo:latest
```

See Makefile for more build options.

## Running the service

```
$ docker run \
  --volume=/var/run/docker.sock:/var/run/docker.sock \
  --env=PLUGIN_SECRET="your shared secret"
  --env=ECR_REGISTRY=xxxxxxxxxx.dkr.ecr.us-east-1.amazonaws.com
```

You should make sure that the service is able to access ECR.

## GCR Integration
The plugin also has the ability to retrieve GCR keys from vault and append them to the response. Just make sure these variables are set:
```
GCR_VAULT_PATH=secret/{container_viewer_json_key}
GCR_REGISTRY_LIST=https://gcr.io,https://us.gcr.io,https://eu.gcr.io,https://asia.gcr.io
VAULT_TOKEN=VAULT_TOKEN
VAULT_ADDR=https://VAULT_HOST:8200
```

## Using the service

To connect your agents to it, specify the following:
```
DRONE_REGISTRY_ENDPOINT="http://your-registry-image:3000"
DRONE_REGISTRY_SECRET="your shared secret"
DRONE_REGISTRY_VERIFY="false"
```
