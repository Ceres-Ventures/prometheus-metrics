FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
WORKDIR $GOPATH/src/github.com/ceres-ventures/prometheus-metrics/
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/metrics cmd/collector/main.go
RUN ls -lah /go/bin/metrics

FROM --platform=linux/x86_64 alpine:3.15.4
RUN apk add --no-cache bash nano
RUN addgroup metrics && adduser -G metrics -D -h /metrics metrics
WORKDIR /metrics
COPY --from=builder /go/bin/metrics /usr/local/bin/metrics
COPY --from=builder /usr/share/zoneinfo/Asia/Almaty /etc/localtime
RUN echo "Asia/Almaty" >  /etc/timezone
COPY --from=builder --chown=metrics:metrics /go/bin/metrics /go/bin/metrics
USER metrics:metrics
ENV BIND_IP=0.0.0.0
ENV BIND_PORT=9292
# Run the hello binary.
RUN ls -lah /usr/local/bin/metrics
CMD ["metrics"]