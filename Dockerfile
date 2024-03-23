FROM golang:1.22-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api-full .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-full . 
EXPOSE 8000
CMD ["./api-full"]