FROM golang:1.21 AS builder

WORKDIR /opt/msgenctl

COPY cmd/ cmd/
COPY internal/ internal/
COPY go.* main.go ./

RUN CGO_ENABLED=0 go build

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /opt/msgenctl/msgenctl /msgenctl

ENTRYPOINT ["/msgenctl"]
