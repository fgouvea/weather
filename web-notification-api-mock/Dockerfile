FROM golang:1.23.5-alpine as builder
RUN apk add build-base

WORKDIR /app
ADD . /app
RUN go build -o /web-notification-api-mock

FROM alpine:3.18
COPY --from=builder /web-notification-api-mock /
EXPOSE 8080
CMD [ "/web-notification-api-mock" ]