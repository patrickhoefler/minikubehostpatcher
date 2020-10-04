# minikubehostpatcher

This is a proof of concept of a solution for the minikube issue [#8439 host.minikube.internal not visible in containers](https://github.com/kubernetes/minikube/issues/8439).

This tool _should_ (only) fix a bug in your local minikube setup. However, please be aware that this is alpha-quality software, so be careful and, when in doubt, check the source code.

## Build

`go build`

## Run

`./minikubehostpatcher`

## License

[MIT](https://github.com/patrickhoefler/minikubehostpatcher/blob/main/LICENSE)
