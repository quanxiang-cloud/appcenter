FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM golang:1.16.6-alpine3.14 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o chaos -mod=vendor -ldflags='-s -w' -installsuffix cgo cmd/chaos/main.go

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

WORKDIR /chaos
COPY --from=builder ./build/chaos ./cmd/

EXPOSE 6666

ENTRYPOINT ["./cmd/chaos", "--kv", "polyurl=http://polyapi:9090/api/v1/polyapi/inner/request/system/poly/permissionInit.p", "--kv", "appcenterurl=http://localhost:80/api/v1/app-center/initCallBack"]