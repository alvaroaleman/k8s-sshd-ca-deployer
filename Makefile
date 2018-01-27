k8s-sshd-ca-deployer: main.go
	CGO_ENABLED=0 go build -ldflags '-s -w' .

bash_unit:
	curl -LO https://raw.githubusercontent.com/pgrange/bash_unit/v1.6.0/bash_unit
	chmod +x bash_unit

.PHONY: test
test: bash_unit k8s-sshd-ca-deployer
	./bash_unit test/integration.sh


fixture-container:
	@if ! docker ps|grep -q k8s-sshd-ca-deployer-fixture-container; then\
		docker run -d \
			-v $$PWD/test/fixtures:/usr/share/nginx/html:ro \
			-p 80:80 \
			--name=k8s-sshd-ca-deployer-fixture-container \
			nginx &>/dev/null; \
	fi

clean:
	-docker rm -f k8s-sshd-ca-deployer-fixture-container


docker_image:
	docker build -t alvaroaleman/k8s-sshd-ca-deployer .
	docker push alvaroaleman/k8s-sshd-ca-deployer
