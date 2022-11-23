# 🪤 k8s-image-replacer ![Go](https://github.com/wuhan005/k8s-image-replacer/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/k8s-image-replacer)](https://goreportcard.com/report/github.com/wuhan005/k8s-image-replacer) [![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph)](https://sourcegraph.com/github.com/wuhan005/k8s-image-replacer)

Replace image when creating pod, aims to speed up image pulling.

## Getting started

### 1. Generate certificate and key

Clone this project. Run `dev/webhook-create-signed-cert.sh` to generate a server certificate/key pair.

```bash
./dev/webhook-create-signed-cert.sh --service k8s-image-replacer --namespace default --secret k8s-image-replacer-tls
```

### 2. Deploy `k8s-image-replacer` into your cluster

```bash
kubectl apply -f dev/k8s-image-replacer.yaml
```

### 3. Create the mutating webhook configuration for `k8s-image-replacer`

Create a mutating webhook with the cluster CA bundle.

```bash
# Replace the ${CA_BUNDLE} placeholder with the actual value.
cat dev/mutating-webhook-template.yaml | ./dev/webhook-patch-ca-bundle.sh > mutating-webhook.yaml

# Create the mutating webhook.
kubectl apply -f mutating-webhook.yaml
```

If you want to set which pod's image to be replaced, check the Kubernetes
document [here](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#matching-requests-objectselector)
.

### (Optional) 4. Build your own image registry proxy.

Follow the [README here](https://github.com/ciiiii/cloudflare-docker-proxy) to build your own image registry proxy with
Cloudflare Workers.

## How it works

1. When a pod is created, the Kubernetes API server will send a request to the
   `k8s-image-replacer` webhook. For more information, check
   the [Kubernetes document](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/).
2. The webhook service will check if the pod's image is in the replacement map in the
   configuration file.
3. If the image is in the replacement map, the webhook will replace the image with the
   replacement image.
4. The pod will be created with the replacement image.
5. The image will be pulled from the image registry proxy.

## Acknowledgements

* [estahn/k8s-image-swapper](https://github.com/estahn/k8s-image-swapper)
* [ciiiii/cloudflare-docker-proxy](https://github.com/ciiiii/cloudflare-docker-proxy)

## License

MIT License
