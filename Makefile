k8s-sshd-ca-deployer: main.go
	CGO_ENABLED=0 go build .

bash_unit:
	curl -LO https://raw.githubusercontent.com/pgrange/bash_unit/v1.6.0/bash_unit
	chmod +x bash_unit

.PHONY: test
test: bash_unit k8s-sshd-ca-deployer
	./bash_unit test/integration.sh
