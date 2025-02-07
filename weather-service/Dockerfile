FROM golang:1.23.5-alpine as builder
RUN apk add build-base

WORKDIR /app
ADD . /app
RUN go build -o /weather-service

FROM alpine:3.18
COPY --from=builder /weather-service /
EXPOSE 8080
CMD [ "/weather-service" ]