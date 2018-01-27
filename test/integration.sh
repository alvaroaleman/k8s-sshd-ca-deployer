#!/usr/bin/env bash

setup_suite(){
  cd ..
  rm -f test/fixtures/empty/*
}

test_init(){
  assert "./k8s-sshd-ca-deployer -ca-url=http://127.0.0.1:80/ca.pem -ca-dest=test/fixtures/empty/cacert -sshd-config-path=test/fixtures/sshd_with_ca"
  assert "echo '62396b67c3423747508951b7ce1c9f26b94f4c5ba8cfab0d8e1dabd15d827b09  test/fixtures/empty/cacert'|sha256sum -c"
}
