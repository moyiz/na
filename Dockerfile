FROM golang:1.21-alpine AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v .

FROM alpine:3.18
COPY --from=builder /usr/src/app/na /usr/local/bin/na
USER 1000
ENTRYPOINT ["na"]