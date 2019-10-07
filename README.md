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
  --env=ECR_REGISTRY_LIST=xxxxxxxxxx.dkr.ecr.us-east-1.amazonaws.com,xxxxxxxxxx.dkr.ecr.us-west-2.amazonaws.com
```

You should make sure that the service is able to access ECR.

## Using the service

To connect your agents to it, specify the following:
```
DRONE_REGISTRY_ENDPOINT="http://your-registry-image:3000"
DRONE_REGISTRY_SECRET="your shared secret"
DRONE_REGISTRY_VERIFY="false"
```
