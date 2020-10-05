package kubectl

import (
	"errors"
	"io"
	"os/exec"
	"strings"
)

// ApplyToKubeSystem applies the provided config to the kube-system namespace
func ApplyToKubeSystem(config string) (output string, err error) {
	applyCmd := exec.Command(
		"kubectl", "--namespace", "kube-system", "apply", "-f", "-",
	)

	stdin, err := applyCmd.StdinPipe()
	if err != nil {
		return
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, config)
	}()

	applyCmdOutput, err := applyCmd.CombinedOutput()
	if err != nil {
		err = errors.New(
			string(applyCmdOutput) + "\n\n" +
				"Tried to apply:\n\n" + config,
		)
	}

	output = string(applyCmdOutput)
	return
}

// GetConfigMap gets the coredns ConfigMap
func GetConfigMap() (configMap string, err error) {
	output, err := exec.Command(
		"kubectl", "--namespace", "kube-system", "get", "configMap/coredns", "-o", "yaml",
	).CombinedOutput()
	if err != nil {
		err = errors.New(string(output))
		return
	}

	configMap = string(output)
	return
}

// GetCurrentContext checks the current kubectl context
func GetCurrentContext() (context string, err error) {
	output, err := exec.Command(
		"kubectl", "config", "current-context",
	).CombinedOutput()
	if err != nil {
		err = errors.New(string(output))
		return
	}

	context = strings.Split(string(output), "\n")[0]
	return
}

// QueryHostIPFromCoreDNS queries the host IP from CoreDNS
func QueryHostIPFromCoreDNS() (hostIP string, err error) {
	output, err := exec.Command(
		"kubectl", "run", "--attach", "--quiet", "--rm", "--restart=Never", "--command",
		"--image=gcr.io/kubernetes-e2e-test-images/dnsutils:1.3@sha256:b31bcf7ef4420ce7108e7fc10b6c00343b21257c945eec94c21598e72a8f2de0", "dnsutils",
		"dig", "+short", "host.minikube.internal",
	).CombinedOutput()
	if err != nil {
		err = errors.New(string(output))
		return
	}

	hostIP = strings.Split(string(output), "\n")[0]
	return
}

// RestartCoreDNSDeployment restarts the deployment to reload the Corefile
func RestartCoreDNSDeployment() (output string, err error) {
	restartDeploymentOutput, err := exec.Command(
		"kubectl", "--namespace", "kube-system", "rollout", "restart", "deployment/coredns",
	).CombinedOutput()
	if err != nil {
		err = errors.New(string(restartDeploymentOutput))
		return
	}

	output = string(restartDeploymentOutput)
	return
}
