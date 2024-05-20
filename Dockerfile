ARG GO_VERSION=1.22
FROM golang:${GO_VERSION} AS builder
WORKDIR /bin
COPY . .
RUN go mod download
RUN go build -o api ./cmd/api

FROM golang:${GO_VERSION} AS build-release-stage
WORKDIR /bin
COPY --from=builder /bin/api .

EXPOSE 8080
ENTRYPOINT [ "/bin/api" ]
