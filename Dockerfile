FROM golang:1.16 AS builder

WORKDIR /opt/msgenctl

COPY cmd/ cmd/
COPY internal/ internal/
COPY go.* main.go ./

RUN CGO_ENABLED=0 go build

FROM scratch

COPY --from=builder /opt/msgenctl/msgenctl /msgenctl

ENTRYPOINT ["/msgenctl"]
