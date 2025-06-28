FROM golang:1.24.4-bookworm AS builder
COPY ./ /go/src/github.com/Vivian-Fan-Jing/mommood
WORKDIR /go/src/github.com/Vivian-Fan-Jing/mommood
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o mommood .

FROM alpine:3.21.3
RUN apk add --no-cache --update --upgrade ca-certificates
RUN adduser -D mommood
USER mommood
COPY --from=builder --chown=mommood:mommood /go/src/github.com/Vivian-Fan-Jing/mommood/mommood /mommood
ENTRYPOINT ["/mommood", "serve"]
