FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o bot

FROM alpine:3.18
WORKDIR /app/
COPY --from=build /app/memories ./memories/
COPY --from=build /app/static ./static/
COPY --from=build /app/templates ./templates/
COPY --from=build /app/bot .
ENV SERVER_PORT=80
CMD ["./bot"]
