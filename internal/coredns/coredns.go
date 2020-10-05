package coredns

const (
	// ExpectedCorefileSnippet is the unpatched original Corefile
	ExpectedCorefileSnippet = `
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        prometheus :9153
`

	// PatchedCorefileSnippetTemplate is the template for the pached Corefile
	PatchedCorefileSnippetTemplate = `
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        hosts {
           %s host.minikube.internal
           fallthrough
        }
        prometheus :9153
`
)
