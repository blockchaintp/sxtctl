FROM golang:1.13-alpine as builder
RUN apk update && \
    apk upgrade && \
    apk add bash
COPY . /tmp/catenasys/sxtctl
WORKDIR /tmp/catenasys/sxtctl
RUN bash ./scripts/build.sh

# ENV CGO_ENABLED=0
# ENV GOOS=linux
# ENV GOARCH=amd64
# ENV GO111MODULE=on
# RUN go build -o /grpcurl \
#     -ldflags "-w -extldflags \"-static\" -X \"main.version=0.0.1\"" \
#     ./cmd/grpcurl

# FROM scratch
# WORKDIR /
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# COPY --from=builder /grpcurl /bin/grpcurl
# ENTRYPOINT ["/bin/grpcurl"]