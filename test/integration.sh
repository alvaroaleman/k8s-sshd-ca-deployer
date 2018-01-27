#!/usr/bin/env bash

setup_suite(){
  cd ..
}

test_init(){
  assert "./k8s-sshd-ca-deployer -ca-url=bla -ca-dest=bla -sshd-config-path=test/fixtures/sshd_with_ca"
}
