# k8s-sidecar-injector

Uses MutatingAdmissionWebhook in Kubernetes to inject sidecars into new deployments at resource creation time.

Forked from [tumblr/k8s-sidecar-injector](https://github.com/tumblr/k8s-sidecar-injector).

# What is this?

It is a small service that runs in each Kubernetes cluster, and listens to the Kubernetes API via webhooks. For each pod creation, the injector gets a (mutating admission) webhook, asking whether or not to allow the pod launch, and if allowed, what changes we would like to make to it. For pods that have special annotations on them (i.e. `injector.tumblr.com/request=logger:v1`), we rewrite the pod configuration to include the containers, volumes, volume mounts, host aliases, init-containers and environment variables defined in the sidecar `logger:v1`'s configuration.

# Deployment

There is a kustomize base provided in [/manifests/base](/manifests/base).

Refer to the [example](/manifests/example) to see what a sample deployment may look like for you.

# Configuration

See [/docs/sidecar-configuration-format.md](/docs/sidecar-configuration-format.md) to get started with setting up your sidecar injector's configurations.

# How it works

1. A pod is created. It has annotation `injector.tumblr.com/request=logger:v1`
2. K8s webhooks out to this service, asking whether to allow this pod creation, and how to mutate it
3. If the pod is annotated with `injector.tumblr.com/status=injected`: Do nothing! Return "allowed" to pod creation
4. Pull the "logger:v1" sidecar config, patch the resource, and return it to k8s
5. Pod will launch in k8s with the modified configuration

A crappy ASCII diagram will help :)

```
                                                                  +-----------------+
     +------------------------------+                             |                 |
     |                              |                             |  Sidecar        |
     |   MutatingAdmissionWebhook   |                             |  configuration  |
     |                              |                             |  files on disk  |
     +------------+-----------------+                             |                 |
                  |                                               +------+----------+
discover injector |                                                      |
endpoints         |                                                      | load from disk
                  |                                                      |
          +-------v--------+    pod launch          +--------------------v-----+
          |                +------------------------>                          |
          |   Kubernetes   |                        |   k8s-sidecar-injector   |
          |   API Server   <------------------------+                          |
          |                |    mutated pod spec    +--------------------------+
          +----------------+
```


# Run

The image is built and published on [Quay](https://quay.io/repository/utilitywarehouse/k8s-sidecar-injector). See [the example](/manifests/example) for how to run this in Kubernetes.

## By hand

```bash
$ ./bin/k8s-sidecar-injector --tls-port=9000 --config-directory=conf/ --tls-cert-file="${TLS_CERT_FILE}" --tls-key-file="${TLS_KEY_FILE}"
```

*NOTE*: this is not a supported method of running in production. You are highly encouraged to deploy this to Kubernetes in [The Supported Way](/manifests/example).

# Hacking

See [hacking.md](/docs/hacking.md)

# License

[Apache 2.0](/LICENSE.txt)

Copyright 2019, Tumblr, Inc.

Copyright (c) 2020 Utility Warehouse Ltd.
