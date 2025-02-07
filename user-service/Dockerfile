FROM golang:1.23.5-alpine as builder
RUN apk add build-base

WORKDIR /app
ADD . /app
RUN go build -o /user-service

FROM alpine:3.18
COPY --from=builder /user-service /
EXPOSE 8080
CMD [ "/user-service" ]