FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/bin/dne-server /usr/local/bin/dne-server

EXPOSE 8080

ENTRYPOINT ["dne-server"]
