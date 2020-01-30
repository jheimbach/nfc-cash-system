FROM golang:alpine as builder
ENV CGO_ENABLED 0
WORKDIR /src
COPY . .

# The -gcflags "all=-N -l" flag helps us get a better debug experience
RUN go build -gcflags "all=-N -l" -o build/server ./cmd/server
RUN apk add --no-cache git
RUN go get github.com/go-delve/delve/cmd/dlv

FROM alpine
WORKDIR run
# 40000 belongs to Delve
EXPOSE 40000
# Allow delve to run on Alpine based containers.
RUN apk add --no-cache libc6-compat

COPY --from=builder src/build/ /run
COPY --from=builder /go/bin/dlv /run

CMD ["./dlv", "--listen=:40000", "--headless=true", "--api-version=2","--accept-multiclient", "exec", "./server"]