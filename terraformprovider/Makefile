TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=chaossearch
NAME=chaossearch
BINARY=terraform-provider-${NAME}
VERSION=0.1.2
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

default: install

build:
	go build -o ${BINARY}

run: install
	cd examples && rm -rf .terraform && terraform init && terraform apply

release:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_${OS_ARCH}


install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
