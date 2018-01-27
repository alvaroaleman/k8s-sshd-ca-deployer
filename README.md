# k8s-sshd-ca-deployer

A simple tool meant to be used inside a Kubernetes DaemonSet to to configure
SSHD to accept keys that were signed by a CA cert from a provided URL.

Used in conjunction with the [vault signed ssh certificates backend](https://www.vaultproject.io/docs/secrets/ssh/signed-ssh-certificates.html)

Checkout the provided `vault_config.sh` to see how Vault has to be configured

# Usage

* Put your vault url into `daemonset.yml`
* `kubectl apply -f daemonset.yml`
