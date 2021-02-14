package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/patrickhoefler/minikubehostpatcher/internal/coredns"
	"github.com/patrickhoefler/minikubehostpatcher/internal/kubectl"
	"github.com/patrickhoefler/minikubehostpatcher/internal/minikube"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	// Remove timestamp from log output
	log.SetFlags(0)

	fmt.Print("Checking if we are in the minikube context ... ")
	currentContext, err := kubectl.GetCurrentContext()
	if err != nil {
		log.Fatal(err)
	}
	if currentContext != "minikube" {
		fmt.Println("❌")
		log.Fatal("Run `kubectl config use-context minikube` and try again.")
	}
	fmt.Println("✅")

	fmt.Println("\nGetting Minikube host IP ...")
	hostIP, err := minikube.GetHostIP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hostIP)

	fmt.Println("\nChecking DNS resolution of host.minikube.internal ...")
	var resolvedHostIP string
	for i := 0; i < 6; i++ {
		resolvedHostIP, err = kubectl.QueryHostIPFromCoreDNS()
		if err == nil {
			break
		}
		fmt.Println("CoreDNS might not be ready yet, trying again ...")
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resolvedHostIP)

	if resolvedHostIP == hostIP {
		fmt.Println(
			"\nhost.minikube.internal resolves correctly, we're all done here 🙂",
		)
		os.Exit(0)
	}

	fmt.Println("DNS resolution of host.minikube.internal is not working yet, let's fix this 😀")

	fmt.Println("\nThis is the CoreDNS Corefile entry we are going to add:")
	patchedCorefileSnippet := fmt.Sprintf(coredns.PatchedCorefileSnippetTemplate, hostIP)

	dmp := diffmatchpatch.New()
	diff := dmp.DiffMain(coredns.ExpectedCorefileSnippet, patchedCorefileSnippet, true)
	fmt.Println(dmp.DiffPrettyText(diff))

	fmt.Println("Getting current Corefile from configMap/coredns ...")
	oldConfigMap, err := kubectl.GetConfigMap()
	if err != nil {
		log.Fatal(err)
	}

	newConfigMap := strings.Replace(oldConfigMap, coredns.ExpectedCorefileSnippet, patchedCorefileSnippet, 1)

	if oldConfigMap == newConfigMap {
		fmt.Println("Error: Corefile was not what we expected.")
		fmt.Println()
		fmt.Println("We were looking for:")
		fmt.Println("...")
		fmt.Println(coredns.ExpectedCorefileSnippet)
		fmt.Println("...")
		fmt.Println()
		fmt.Println("We received:")
		fmt.Println("\n" + oldConfigMap)
		os.Exit(1)
	}

	fmt.Println("\nPatching Corefile ...")
	dmp = diffmatchpatch.New()
	diff = dmp.DiffMain(oldConfigMap, newConfigMap, true)
	fmt.Println(dmp.DiffPrettyText(diff))

	fmt.Println("Replacing patched Corefile ...")
	output, err := kubectl.ReplaceInKubeSystem(newConfigMap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

	fmt.Println("Restarting coredns deployment to pick up the change ...")
	output, err = kubectl.RestartCoreDNSDeployment()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(output)

	fmt.Println("\nChecking DNS resolution of host.minikube.internal ...")
	for i := 0; i < 3; i++ {
		resolvedHostIP, err = kubectl.QueryHostIPFromCoreDNS()
		if err == nil {
			break
		}
		fmt.Println("CoreDNS might not be ready yet, trying again ...")
		time.Sleep(10)
	}
	if err != nil {
		log.Fatal(err)
	}

	if resolvedHostIP != hostIP {
		fmt.Println("Sorry, somehow that didn't quite work 😕")
		fmt.Println(fmt.Sprintf("We expected %s, but received %s", hostIP, resolvedHostIP))
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("host.minikube.internal now resolves to %s 🙂", hostIP))
}
