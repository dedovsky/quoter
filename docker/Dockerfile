FROM golang:1.24-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash gcc musl-dev

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY ./ ./

RUN go build -o quoter ./cmd

FROM alpine:latest AS runner

COPY --from=builder /usr/local/src/quoter /

ENTRYPOINT ["./quoter"]