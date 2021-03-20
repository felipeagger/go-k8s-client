# go-k8s-client
Golang managing resources in kubernetes

List Nodes, Pods.. Create & Delete Pods, Jobs...

# Requirements

Add k8s config file in root directory as "kube_config.yml":

```
apiVersion: v1
kind: Config
clusters:
- cluster:
    api-version: v1
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0t...tLS0K
    server: "https://192.168.100.50:6443"
  name: "local"
contexts:
- context:
    cluster: "local"
    user: "kube-admin-local"
  name: "local"
current-context: "local"
users:
- name: "kube-admin-local"
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQlNREV3TWpReE5ERXpNV....UtLS0tLQo=
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJ......==
```

# Running:

```
go run main.go
```
