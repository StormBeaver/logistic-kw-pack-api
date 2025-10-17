# Builder
FROM golang:1.24.1-alpine AS builder

ARG GITHUB_PATH=github.com/StormBeaver/logistic-pack-api

WORKDIR /home/${GITHUB_PATH}

RUN apk add --update make git protoc protobuf protobuf-dev curl
COPY Makefile Makefile
RUN make deps-go
COPY . .
RUN make build-go

# gRPC Server

FROM alpine:latest AS server

ARG GITHUB_PATH=github.com/StormBeaver/logistic-pack-api

LABEL org.opencontainers.image.source=https://${GITHUB_PATH}

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/grpc-server .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .
COPY --from=builder /home/${GITHUB_PATH}/migrations/ ./migrations

RUN chown root:root grpc-server

EXPOSE 8080
EXPOSE 8082
EXPOSE 9100
EXPOSE 8000

CMD ["./grpc-server"]
