# minikubehostpatcher

[![Go Report Card](https://goreportcard.com/badge/github.com/patrickhoefler/minikubehostpatcher)](https://goreportcard.com/report/github.com/patrickhoefler/minikubehostpatcher)
[![Maintainability](https://api.codeclimate.com/v1/badges/af9c56e5eb950771cc56/maintainability)](https://codeclimate.com/github/patrickhoefler/minikubehostpatcher/maintainability)

This is a proof of concept of a solution for the Minikube issue [#8439 host.minikube.internal not visible in containers](https://github.com/kubernetes/minikube/issues/8439).

According to the Minikube documentation page [Host access](https://minikube.sigs.k8s.io/docs/handbook/host-access/), it should be possible to get the host machine IP address using `host.minikube.internal` from inside pods. Unfortunately, this mechanism is currently broken. `minikubehostpatcher` remedies this situation.

Please be aware that this tool is alpha-quality software, so be careful and, when in doubt, check the source code.

This prototype is continously tested on macOS 10.15.6 with Minikube 1.13.0, Virtualbox 6.1.14 and the following Kubernetes versions:

- 1.19.2
- 1.18.9
- 1.17.12
- 1.16.15
- 1.15.12

In addition, it has been manually tested on macOS 10.15.6 with Minikube 1.13.1, Kubernetes 1.19.2 and the following drivers:

- Docker (19.03.12)
- hyperkit (0.20200224-44-gb54460)

Currently, it is not compatible with Kubernetes 1.14 and earlier.

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

This is the Corefile entry we are going to add:

           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        hosts {
           192.168.65.2 host.minikube.internal
           fallthrough
        }
        prometheus :9153

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
  creationTimestamp: "2020-10-05T12:30:56Z"
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:data:
        .: {}
        f:Corefile: {}
    manager: kubeadm
    operation: Update
    time: "2020-10-05T12:30:56Z"
  name: coredns
  namespace: kube-system
  resourceVersion: "212"
  selfLink: /api/v1/namespaces/kube-system/configmaps/coredns
  uid: a5a7e3ba-4315-4364-88d0-070311e92aa5

Replacing patched Corefile ...
configmap/coredns replaced

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
