FROM golang:alpine3.13 AS go-builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=go-builder /app/main .
CMD ["/app/main"]
