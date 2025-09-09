FROM node:22-alpine AS builder-ui

RUN apk --update add make
WORKDIR /build
ADD . .
RUN make build-ui

FROM golang:1.24-alpine AS builder

RUN apk --update add make
WORKDIR /build
ADD . .
COPY --from=builder-ui /build/frontend/ui/dist ./frontend/ui/dist
RUN make build-docker && \
    cd bin/release && \
    mv *_docker rtcserver

FROM alpine:3.15.0

WORKDIR /rtcserver
COPY examples/config.yaml config.yaml
COPY --from=builder /build/bin/release .
CMD ["./rtcserver"]