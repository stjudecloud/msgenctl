FROM golang:1.17 AS builder

WORKDIR /opt/msgenctl

COPY cmd/ cmd/
COPY internal/ internal/
COPY go.* main.go ./

RUN CGO_ENABLED=0 go build

FROM alpine:3.14.2

ENV PATH=/opt/msgenctl/bin:$PATH
COPY --from=builder /opt/msgenctl/msgenctl /opt/msgenctl/bin/msgenctl

# Added for workflow runners.
RUN apk add --no-cache bash

ENTRYPOINT ["/opt/msgenctl/bin/msgenctl"]
