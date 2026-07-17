FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o comment-service .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/comment-service .
EXPOSE 50051
CMD ["./comment-service"]
