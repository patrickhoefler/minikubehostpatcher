package minikube

import (
	"errors"
	"os/exec"
	"regexp"
)

// GetHostIP gets the host IP from Minikube
func GetHostIP() (hostIP string, err error) {
	getMinikubeHostOutput, err := exec.Command(
		"minikube", "ssh", "grep", "host.minikube.internal", "/etc/hosts",
	).CombinedOutput()
	if err != nil {
		err = errors.New(string(getMinikubeHostOutput))
		return
	}

	hostIP = regexp.MustCompile(`^\S+`).FindString(string(getMinikubeHostOutput))
	return
}
