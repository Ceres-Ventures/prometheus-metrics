FROM --platform=linux/x86_64 golang:alpine AS builder
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
WORKDIR $GOPATH/src/github.com/ceres-ventures/prometheus-metrics/
COPY . .
RUN go mod download
#RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/metrics cmd/collector/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -installsuffix 'static' -o /app cmd/collector/main.go

FROM --platform=linux/x86_64 scratch AS final
LABEL maintainer="gbaeke"
COPY --from=builder  /app /metrics
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/metrics"]