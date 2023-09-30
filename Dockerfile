FROM golang:1.19 AS build

WORKDIR /app

COPY . .

RUN export CGO_ENABLED=0 && go build -ldflags="-s -w" -o main main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 8088

CMD ["./main"]