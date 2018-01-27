#!/usr/bin/env bash

cd $(mktemp -d)

# Check out the vault doc on how to set up a vault server properly
# Keep in mind that your Vault should be highly available, you can't access your
# servers anymore if its down!
vault server -dev &
sleep 10s

export VAULT_ADDR=http://127.0.0.1:8200

# Configure the backend
vault mount -path=ssh-client-signer ssh
vault write ssh-client-signer/config/ca generate_signing_key=true

# Create the backend role
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
  "max_ttl": "30m"
}
EOH

# You can now use the provided tool to configure your sshds like this:
./k8s-sshd-ca-deployer \
  -ca-url=http://127.0.0.1:8200/v1/ssh-client-signer/public_key \
  -ca-dest=/etc/sshd_cert \
  -sshd-config-path=/etc/ssh/sshd_config \
  -restart-command="systemctl restart sshd"

# Create a keypair to use
yes|ssh-keygen -f $PWD/id_rsa_vault -b 4096

# Request a certificate, in reality you
# are going to have this as an alias in your ~/.bashrc
vault write -field=signed_key \
  ssh-client-signer/sign/centos \
  public_key=@$PWD/id_rsa_vault.pub \
  valid_principals=root,centos,core,ubuntu \
  allow_user_key_ids=true \
    > $PWD/id_rsa_vault-cert.pub

# You can now test login
ssh -i $PWD/id_rsa_vault root@localhost "echo awesome stuff is awesome!"

# Kill the vault dev instance
ps xf|grep vault|grep -v grep|awk '{ print $1 }'|xargs kill
