## Creating controller scaffolding
```bash
kubebuilder init --domain junaidk.io --repo github.com/junaidk/image-backup-controller
```

## Makefile Options 

### Generating manifests
`make manifests`

### Building docker image
`make docker-build`

### Deploying to cluster
`make deploy`

### Removing controller
`make undeploy`

## Running tests
`make test`

## Image Copy
https://github.com/containers/image

This library is used to copy image from source to destination.

## Running the operator

Build the image using `make docker-build IMG=<some-registry>/<project-name>:tag`

Push the image to the registry using `make docker-push IMG=<some-registry>/<project-name>:tag`

Create secret with registry username and password.

```bash
kubectl create secret generic registry-creds \
  --from-literal=username=devuser \
  --from-literal=password='S!B\*d$zDsb='
```

Update environment variables in config/manager/manager.yaml according to your destination registry.

And update the secret name in config/manager/manager.yaml

Additional namespaces can be added to env IGNORE_NAMESPACES in config/manager/manager.yaml. These will be ignored by controller in addtion to `kube-system`

Deploy the controller to the cluster using `make deploy IMG=<some-registry>/<project-name>:tag`

## Improvements

- Include Init container in image update process.
- Make image copy concurrent in case of multiple images in one deployment.

## asciinema Recording

https://asciinema.org/a/p3HpSVxuKk3duf3zCjGQStsLu
