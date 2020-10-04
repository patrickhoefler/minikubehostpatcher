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

	fmt.Println("\nChecking CoreDNS resolution of host.minikube.internal ...")
	resolvedHostIP, err := kubectl.QueryHostIPFromCoreDNS()
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

	fmt.Println("CoreDNS resolution of host.minikube.internal is not working yet, let's fix this 😀")

	fmt.Println("\nThis is the patch we are going to apply:")
	patchedData := fmt.Sprintf(coredns.PatchedDataTemplate, hostIP)

	dmp := diffmatchpatch.New()
	diff := dmp.DiffMain(coredns.ExpectedData, patchedData, true)
	fmt.Println(dmp.DiffPrettyText(diff))

	fmt.Println("Getting current Corefile from configMap/coredns ...")
	oldConfigMap, err := kubectl.GetConfigMap()
	if err != nil {
		log.Fatal(err)
	}

	newConfigMap := strings.Replace(oldConfigMap, coredns.ExpectedData, patchedData, 1)

	if oldConfigMap == newConfigMap {
		fmt.Println("Error: Corefile was not what we expected.")
		fmt.Println()
		fmt.Println("We were looking for:")
		fmt.Println(coredns.ExpectedData)
		fmt.Println()
		fmt.Println("We received:")
		fmt.Println("\n" + oldConfigMap)
		os.Exit(1)
	}

	fmt.Println("\nPatching Corefile ...")
	dmp = diffmatchpatch.New()
	diff = dmp.DiffMain(oldConfigMap, newConfigMap, true)
	fmt.Println(dmp.DiffPrettyText(diff))

	fmt.Println("Applying patched Corefile ...")
	output, err := kubectl.ApplyToKubeSystem(newConfigMap)
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

	fmt.Println("\nChecking CoreDNS resolution of host.minikube.internal ...")
	for i := 0; i < 3; i++ {
		resolvedHostIP, err = kubectl.QueryHostIPFromCoreDNS()
		if err == nil {
			break
		}
		if strings.Index(err.Error(), "connection timed out") >= 0 {
			fmt.Println("CoreDNS not ready yet, trying again ...")
			time.Sleep(1)
		} else {
			log.Fatal(err)
		}
	}

	if resolvedHostIP != hostIP {
		fmt.Println("Sorry, somehow that didn't quite work 😕")
		fmt.Println(fmt.Sprintf("We expected %s, but received %s", hostIP, resolvedHostIP))
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("host.minikube.internal now resolves to %s 🙂", hostIP))
}
