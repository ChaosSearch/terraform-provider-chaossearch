#!/bin/sh
export TF_LOG=DEBUG && export TF_LOG_PATH="terraform.txt"
export TF_VAR_admin_password=$CLUSTER_ADMIN_PASSWORD TF_VAR_password=$CLUSTER_USER_PASSWORD TF_VAR_secret_access_key=$CLUSTER_SECRET_ACCESS_KEY
shopt -s dotglob

find * -prune -type d | while IFS= read -r d; do
  echo "$PWD/$d"
  cd "$PWD/$d"

  rm .terraform.lock.hcl

  rm terraform.tfstate
  rm terraform.tfstate.backup
  rm terraform.txt

  for file in *.tf; do
    terraform init
    echo "-------------------------------Start Execution of $file -------------------------------";
    terraform apply -auto-approve -var-file ../../terraform-dev.tfvars -compact-warnings
    terraform destroy -auto-approve -var-file ../../terraform-dev.tfvars
    echo "-------------------------------End Execution of $file   -------------------------------";
  done
  cd ..
done