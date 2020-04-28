FROM golang:1.14.2-alpine AS build

# Set workdir
WORKDIR /work

# Add dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 go build ./cmd/vault-plugin-secrets-minio

FROM vault:1.4.0
COPY --from=build /work/vault-plugin-secrets-minio /plugins/vault-plugin-secrets-minio
