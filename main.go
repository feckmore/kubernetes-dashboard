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

	dashboardPodName := bash("\nLooking for previous deployment...", "kubectl get pods -n kubernetes-dashboard | awk '/kubernetes-dashboard/ {print $1}'", true, true)
	dashboardPodName = strings.TrimSpace(dashboardPodName)
	if len(dashboardPodName) == 0 {
		fmt.Println("Previous deployment not found.")
		bash("\nDeploying Dashboard...", "kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta4/aio/deploy/recommended.yaml", true, true)
	} else {
		fmt.Println("Previous deployment found.")
	}

	userSecretsName := bash("\nLooking for existing admin user secrets...", "kubectl get secret -n kubernetes-dashboard | awk '/admin-user/ {print $1}'", true, false)
	userSecretsName = strings.TrimSpace(userSecretsName)
	if len(userSecretsName) == 0 {
		fmt.Println("Existing admin user secrets not found.")

		bash("\nCreating Admin User...", "kubectl apply -f ./adminuser.serviceaccount.yml", true, true)
		userSecretsName = bash("\nParsing admin user secrets name...", "kubectl get secret -n kubernetes-dashboard | awk '/admin-user/ {print $1}'", true, true)
		userSecretsName = strings.TrimSpace(userSecretsName)
	} else {
		fmt.Println("Existing admin user secrets found: ", userSecretsName)
	}

	token := bash("\nParsing user token...", fmt.Sprintf("kubectl describe secret %s -n kubernetes-dashboard | awk '/token:/ {print $2}'", userSecretsName), true, false)

	// opening browser link before actually available via proxy. This allows proxy step, which does not complete, to be the last step
	bash("\nOpening Dashboard in Browser", "open http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/", true, true)

	fmt.Printf("%s\n%s\n", "Token to paste into Dashboard:", token)

	bash("\nProxy to Dashboard...", "kubectl proxy", true, true)
}

// bash executes a command with bash, optionally printing out a description, the command itself and stdout/stderr to the console
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
