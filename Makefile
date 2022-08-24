TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=chaossearch
NAME=chaossearch
BINARY=terraform-provider-${NAME}
VERSION=1.0.6
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

default: install

lintProvider:
	golangci-lint run -c .golangci.yml

build: lintProvider
	 go build -o ${BINARY}

runEx: install
	cd provider/examples && terraform init && terraform apply --auto-approve

runExDeleteTF: install
	cd provider/examples && rm -rf .terraform .terraform* *.tfstate* && terraform init && terraform apply --auto-approve

release:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_${OS_ARCH}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

initTestInfra:
	cd provider/tests/test-infra && terraform init && terraform apply

testAcc: install
	TF_ACC=1 go test -v ./provider/tests/