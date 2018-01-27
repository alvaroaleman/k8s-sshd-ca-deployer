#!/usr/bin/env bash

cd $(mktemp -d)

vault server -dev &
sleep 10s

export VAULT_ADDR=http://127.0.0.1:8200

vault mount -path=ssh-client-signer ssh
vault write ssh-client-signer/config/ca generate_signing_key=true

vault write ssh-client-signer/roles/centos -<<"EOH"
{
  "allow_user_certificates": true,
  "allowed_users": "*",
  "default_extensions": [
    {
      "permit-pty": "",
      "permit-port-forwarding": "",
      "permit-agent-forwarding": ""
    }
  ],
  "key_type": "ca",
  "default_user": "centos",
  "max_ttl": "10s"
}
EOH

yes|ssh-keygen -f $PWD/id_rsa_vault -b 4096

vault write -field=signed_key \
  ssh-client-signer/sign/centos \
  public_key=@$PWD/id_rsa_vault.pub \
  valid_principals=root,centos,core \
  allow_user_key_ids=true \
    > $PWD/id_rsa_vault-cert.pub

ssh-add $PWD/id_rsa_vault
sleep 30s
ssh-add -d $PWD/id_rsa_vault

ps xf|grep vault|grep -v grep|awk '{ print $1 }'|xargs kill
