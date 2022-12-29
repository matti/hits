FROM golang:1.19-alpine3.17 as builder

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) go build -o /hits

FROM alpine:3.17

RUN apk add --no-cache curl
COPY --from=builder /hits /usr/local/bin

ENTRYPOINT [ "/usr/local/bin/hits" ]
