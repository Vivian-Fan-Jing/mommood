FROM golang:1.24.4-bookworm AS builder
COPY ./ /go/src/github.com/Vivian-Fan-Jing/mommood
WORKDIR /go/src/github.com/Vivian-Fan-Jing/mommood
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o mommood .

FROM alpine:3.21.3
RUN apk add --no-cache --update --upgrade ca-certificates
RUN mkdir -p /pb_data
ADD ./pb_data /pb_data
RUN ls -lrt /pb_data
RUN adduser -D mommood
USER mommood
COPY --from=builder --chown=mommood:mommood /go/src/github.com/Vivian-Fan-Jing/mommood/mommood /mommood
VOLUME /pb_data
EXPOSE 8090
ENTRYPOINT ["/mommood", "serve", "--http=0.0.0.0:8090"]
