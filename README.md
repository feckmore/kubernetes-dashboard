# Kubernetes Dashboard Scripts

Golang command line app, that uses `bash` calls to bring up the Kubernetes Web UI Dashboard.

Taken from:
https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/

https://github.com/kubernetes/dashboard/blob/master/docs/user/access-control/creating-sample-user.md

Assumptions:

- MacOS
- Using Docker for Desktop locally for Kubernetes
- `kubectl` is installed
- `bash` & `awk` available

This could just use a bash script, but where's the fun in that?
