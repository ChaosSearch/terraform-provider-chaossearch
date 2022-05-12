TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=chaossearch
NAME=chaossearch
BINARY=terraform-provider-${NAME}
VERSION=0.1.1
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

default: install

lintProvider:
	golangci-lint run -c .golangci.yml

build: lintProvider
	 go build -o ${BINARY}

runExample: install
	cd provider/examples && rm -rf .terraform && terraform init && terraform apply -var-file ../examples/terraform-dev.tfvars

release:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_${OS_ARCH}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

initTestInfra:
	cd provider/tests/test-infra && terraform init && terraform apply

testAcc: install
	TF_ACC=1 go test -v ./provider/tests/