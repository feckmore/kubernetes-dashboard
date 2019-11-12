package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const CONTEXT = "docker-desktop"

func main() {
	// setting the context to local docker cluster, so dashboard will be created there
	bash("Setting Context...", fmt.Sprintf("%s %s", "kubectl config set-context", CONTEXT), true, true)

	bash("Deploying Dashboard...", "kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta4/aio/deploy/recommended.yaml", true, true)

	bash("Creating Admin User...", "kubectl apply -f ./adminuser.serviceaccount.yml", true, true)

	tokenPodName := bash("Finding User Token Pod Name...", "kubectl get secret -n kubernetes-dashboard | awk '/admin-user/ {print $1}'", true, true)
	token := bash("Getting User Token...", fmt.Sprintf("kubectl describe secret %s -n kubernetes-dashboard | awk '/token:/ {print $2}'", tokenPodName), true, false)

	// opening browser link before actually available via proxy. This allows proxy step, which does not complete, to be the last step
	bash("Opening Dashboard in Browser", "open http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/", true, true)

	fmt.Printf("%s\n%s\n\n", "Token to paste into Dashboard:", token)

	bash("Proxy to Dashboard...", "kubectl proxy", true, true)
}

func bash(description, commandString string, printIn, printOut bool) (out string) {
	if len(description) > 0 {
		fmt.Println(description)
	}

	if printIn {
		fmt.Println(commandString)
	}
	cmd := exec.Command("bash", "-c", commandString)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	out = strings.TrimSpace(string(stdoutStderr))

	if printOut {
		fmt.Println(out)
	}

	return out
}
