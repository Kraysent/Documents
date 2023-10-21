FROM golang:1.20 as builder
WORKDIR /build/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./documents

FROM ubuntu:latest as runner
WORKDIR /app/
COPY --from=builder /build/documents .
ADD "https://storage.yandexcloud.net/cloud-certs/CA.pem" "/.ssl/root.crt"
CMD ["/app/documents"]
