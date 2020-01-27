FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN go build -o build/server cmd/main.go

FROM alpine
WORKDIR run
COPY --chown=0:0 --from=builder /src/build/ /run

RUN apk add --update bash curl && rm -rf /var/cache/apk/*

ENTRYPOINT ["./server"]