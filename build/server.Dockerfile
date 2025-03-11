FROM golang:1.22-alpine AS proto-builder
RUN apk update && apk add protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /usr/src/app

COPY proto proto
RUN mkdir -p pkg/proto
RUN protoc --proto_path=proto --go_out=pkg/proto --go_opt=paths=source_relative --go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative chat.proto

FROM golang:1.22-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/server/ cmd/server/
COPY internal/ internal/
COPY pkg/ pkg/
COPY --from=proto-builder /usr/src/app/pkg/proto pkg/proto

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/main.go

FROM scratch

COPY --from=builder /usr/src/app/server /server

ENTRYPOINT ["/server"]
