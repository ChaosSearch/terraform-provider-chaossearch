#!/bin/sh

rm -rf .terraform
rm -rf .terraform.lock.hcl
rm -rf terraform.txt
rm -rf terraform.tfstate

echo "Deleted chache files..."

echo "---------------------------------------------------"
ls -a
echo "---------------------------------------------------"

cd ..
make install
cd examples

terraform init --upgrade
terraform apply
