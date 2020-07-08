##########################################
FROM golang:1.14.4-alpine3.12 as builder

WORKDIR /go/src/app

COPY . .

WORKDIR /go/src/app/cmd/kitops/

RUN go build -v

##########################################
FROM alpine:3.12.0

COPY scripts /scripts

COPY --from=builder /go/src/app/cmd/kitops/kitops /

ENTRYPOINT [ "/scripts/entrypoint.sh" ]

EXPOSE 8080/tcp
