# Build the manager binary
FROM golang:1.24 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/operator/main.go cmd/operator/main.go
COPY pkg/ pkg/
COPY must-gather/ must-gather/

# Build
RUN GOOS=linux GOARCH=amd64 go build -a -tags netgo,osusergo -o manager cmd/operator/main.go

# Use ubi. The gather script requires bash
FROM registry.access.redhat.com/ubi8/ubi-micro
WORKDIR /

COPY --from=builder /workspace/manager .
COPY --from=builder /workspace/must-gather/collection-scripts/* .
USER 65532:65532

ENTRYPOINT ["/manager"]
