FROM scratch
WORKDIR /

COPY ./target /target
ENTRYPOINT ["/target/sxtctl-linux-amd64"]
