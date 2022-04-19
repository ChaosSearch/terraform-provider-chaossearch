TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=chaossearch
NAME=chaossearch
BINARY=terraform-provider-${NAME}
VERSION=0.1.1
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

default: install

lint_provider:
	golangci-lint run -c .golangci.yml

build: lint_provider
	 go build -o ${BINARY}

run: install
	cd provider/examples && rm -rf .terraform && terraform init && terraform apply -var-file ../examples/terraform-dev.tfvars

release:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_${OS_ARCH}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

testall:
	cd examples/test && ./testall.sh