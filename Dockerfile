FROM golang:latest AS build

WORKDIR /app

COPY . .

RUN export CGO_ENABLED=0 && go build -ldflags="-s -w" -o random_image random_image.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/random_image .

EXPOSE 8088

CMD ["./random_image"]