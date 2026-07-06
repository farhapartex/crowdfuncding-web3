FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o crowdfunding-backend .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/crowdfunding-backend .
EXPOSE 8080
CMD ["./crowdfunding-backend"]
