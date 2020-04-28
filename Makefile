PROJECT		= "github.com/statcan/vault-plugin-secrets-minio"
GOFILES		= $(shell find . -name "*.go")

default: vault-plugin-secrets-minio

vault-plugin-secrets-minio: $(GOFILES)
	go build ./cmd/vault-plugin-secrets-minio

clean:
	rm -f vault-plugin-secrets-minio

deps:
	go get ./...

.PHONY: default clean deps
