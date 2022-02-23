#!/bin/sh

shopt -s dotglob
find * -prune -type d | while IFS= read -r d; do
  echo "$PWD/$d"
  cd "$PWD/$d"

  rm .terraform.lock.hcl
  for file in *; do

    terraform init
    echo "-------------------------------Start Execution of $file -------------------------------";
    terraform apply -auto-approve
    terraform destroy -auto-approve
    echo "-------------------------------End Execution of $file   -------------------------------";
  done
  cd ..
done