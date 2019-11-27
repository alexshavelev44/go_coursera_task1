# docker build -t mailgo_hw1 .
FROM golang:1.9.2-alpine
COPY . .
RUN go test -v