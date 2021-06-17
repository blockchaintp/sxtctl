FROM golang:1.16-alpine as builder
RUN apk update && \
    apk upgrade && \
    apk add \
      bash \
      gcc
COPY . /tmp/catenasys/sxtctl
WORKDIR /tmp/catenasys/sxtctl
RUN bash ./scripts/build

FROM scratch
WORKDIR /
COPY --from=builder /tmp/catenasys/sxtctl/target /target
ENTRYPOINT ["/target/sxtctl-linux-amd64"]
