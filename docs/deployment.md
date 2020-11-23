# Deployment

Example Kubernetes manifests are provided in [/examples/kubernetes](/examples/kubernetes). You are expected to tailor these to your needs. Specifically, you will need to:

1. Generate TLS certs [/docs/tls.md](/docs/tls.md) and update [/examples/kubernetes/mutating-webhook-configuration.yaml](/examples/kubernetes/mutating-webhook-configuration.yaml) with the `caBundle`
2. Update [/examples/kubernetes/deployment.yaml](/examples/kubernetes/deployment.yaml) with the appropriate version you want to deploy
3. Specify whatever flags you want in the deployment.yaml
4. Create a kubernetes secret from the certificates that you generated as a part of [/docs/tls.md](/docs/tls.md).
```
kubectl create secret generic k8s-sidecar-injector --from-file=examples/tls/${DEPLOYMENT}/${CLUSTER}/sidecar-injector.crt --from-file=examples/tls/${DEPLOYMENT}/${CLUSTER}/sidecar-injector.key --namespace=kube-system
```
5. Create a ConfigMap containing some config files and mount it in the container so the injector has some sidecars to inject :) [/docs/sidecar-configuration-format.md](/docs/sidecar-configuration-format.md)

Once you hack the example Kubernetes manifests to work for your deployment, deploy them to your cluster. The list of manifests you should deploy are below:

* [clusterrole.yaml](/examples/kubernetes/clusterrole.yaml)
* [clusterrolebinding.yaml](/examples/kubernetes/clusterrolebinding.yaml)
* [service-monitor.yaml](/examples/kubernetes/service-monitor.yaml)
* [serviceaccount.yaml](/examples/kubernetes/serviceaccount.yaml)
* [service.yaml](/examples/kubernetes/service.yaml)
* [deployment.yaml](/examples/kubernetes/deployment.yaml)
* [mutating-webhook-configuration.yaml](/examples/kubernetes/mutating-webhook-configuration.yaml)

A sample ConfigMap is included to test injections at [/examples/kubernetes/configmap-sidecar-test.yaml](/examples/kubernetes/configmap-sidecar-test.yaml).

Now, you are ready to create your first pod that asks for an injection:

```bash
$ kubectl create -f examples/kubernetes/debug-pod.yaml
pod/debian-debug created
```

Verify its up and running; note the `injector.tumblr.com/status: injected` label, indicating the pod had its sidecar added successfully, as well as the added environment variables, and additional `sidecar-nginx` container!

```bash
$ kubectl describe -f debug-pod.yaml
Name:         debian-debug
Namespace:    default
...
Annotations:  injector.tumblr.com/status: injected
Status:       Running
IP:           10.246.248.115
Containers:
  debian-debug:
    Image:         debian:jessie
    Command:
      sleep
      3600
    State:          Running
      Started:      Mon, 19 Nov 2018 11:28:36 -0500
    Ready:          True
    Restart Count:  0
    Environment:
      HELLO:  world
  sidecar-nginx:
    Image:          nginx:1.12.2
    Port:           80/TCP
    Host Port:      0/TCP
    State:          Running
      Started:      Mon, 19 Nov 2018 11:28:40 -0500
    Ready:          True
    Environment:
      HELLO:  world
    Mounts:   <none>
...
```


