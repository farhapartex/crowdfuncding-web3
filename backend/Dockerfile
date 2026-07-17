FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.work go.work.sum ./
COPY backend ./backend
COPY comment ./comment
WORKDIR /app/backend
RUN CGO_ENABLED=0 go build -o crowdfunding-backend .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/backend/crowdfunding-backend .
EXPOSE 8080
CMD ["./crowdfunding-backend"]
