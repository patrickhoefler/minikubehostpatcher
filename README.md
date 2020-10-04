# minikubehostpatcher

[![Go Report Card](https://goreportcard.com/badge/github.com/patrickhoefler/minikubehostpatcher)](https://goreportcard.com/report/github.com/patrickhoefler/minikubehostpatcher)
[![Maintainability](https://api.codeclimate.com/v1/badges/af9c56e5eb950771cc56/maintainability)](https://codeclimate.com/github/patrickhoefler/minikubehostpatcher/maintainability)

This is a proof of concept of a solution for the minikube issue [#8439 host.minikube.internal not visible in containers](https://github.com/kubernetes/minikube/issues/8439).

This tool _should_ (only) fix a bug in your local minikube setup. However, please be aware that this is alpha-quality software, so be careful and, when in doubt, check the source code.

So far, this prototype has been successfully tested with minikube v1.13.1 on:

- macOS (10.15.6)
  - Docker (19.03.12)
  - hyperkit (0.20200224-44-gb54460)
  - VirtualBox (6.1.14)

_In theory_, this approach should work everywhere, since it only adds a host line to the CoreDNS Corefile.

## Build

`go build`

## Run

`./minikubehostpatcher`

## Output

```text
❯ minikube start
😄  minikube v1.13.1 on Darwin 10.15.6
✨  Automatically selected the docker driver
👍  Starting control plane node minikube in cluster minikube
🔥  Creating docker container (CPUs=2, Memory=1990MB) ...
🐳  Preparing Kubernetes v1.19.2 on Docker 19.03.8 ...
🔎  Verifying Kubernetes components...
🌟  Enabled addons: default-storageclass, storage-provisioner
🏄  Done! kubectl is now configured to use "minikube" by default


❯ ./minikubehostpatcher
Checking if we are in the minikube context ... ✅

Getting Minikube host IP ...
192.168.65.2

Checking CoreDNS resolution of host.minikube.internal ...

CoreDNS resolution of host.minikube.internal is not working yet, let's fix this 😀

This is the patch we are going to apply:

  Corefile: |
    .:53 {
        errors
        health {
           lameduck 5s
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        hosts {
           192.168.65.2 host.minikube.internal
           fallthrough
        }
        prometheus :9153
        forward . /etc/resolv.conf {
           max_concurrent 1000
        }
        cache 30
        loop
        reload
        loadbalance
    }

Getting current Corefile from configMap/coredns ...

Patching Corefile ...
apiVersion: v1
data:
  Corefile: |
    .:53 {
        errors
        health {
           lameduck 5s
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        hosts {
           192.168.65.2 host.minikube.internal
           fallthrough
        }
        prometheus :9153
        forward . /etc/resolv.conf {
           max_concurrent 1000
        }
        cache 30
        loop
        reload
        loadbalance
    }
kind: ConfigMap
metadata:
  creationTimestamp: "2020-10-04T17:32:29Z"
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:data:
        .: {}
        f:Corefile: {}
    manager: kubeadm
    operation: Update
    time: "2020-10-04T20:08:26Z"
  name: coredns
  namespace: kube-system
  resourceVersion: "4098"
  selfLink: /api/v1/namespaces/kube-system/configmaps/coredns
  uid: 8589b276-c0b3-41c1-9e02-20bb65d1e23d

Applying patched Corefile ...
Warning: kubectl apply should be used on resource created by either kubectl create --save-config or kubectl apply
configmap/coredns configured

Restarting coredns deployment to pick up the change ...
deployment.apps/coredns restarted

Checking CoreDNS resolution of host.minikube.internal ...
host.minikube.internal now resolves to 192.168.65.2 🙂


❯ ./minikubehostpatcher
Checking if we are in the minikube context ... ✅

Getting Minikube host IP ...
192.168.65.2

Checking CoreDNS resolution of host.minikube.internal ...
192.168.65.2

host.minikube.internal resolves correctly, we're all done here 🙂
```

## License

[MIT](https://github.com/patrickhoefler/minikubehostpatcher/blob/main/LICENSE)
