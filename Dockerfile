# build stage
FROM golang:1.13-alpine AS builder
RUN apk --no-cache add git
RUN GO111MODULE=on CGO_ENABLED=0 go build -ldflags="-s -w" github.com/paymentdata/releaseforms

# final stage
FROM scratch
COPY --from=builder /go/releaseforms /
ENTRYPOINT ["/releaseforms"]