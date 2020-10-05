package coredns

const (
	// ExpectedCorefileSnippet is the unpatched original Corefile
	ExpectedCorefileSnippet = `
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        prometheus :9153
`

	// PatchedCorefileSnippetTemplate is the template for the pached Corefile
	PatchedCorefileSnippetTemplate = `
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
