FROM golang:1.21 as builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o ./main

WORKDIR /app


RUN mkdir publish
RUN cp main publish
RUN cp -r config publish

FROM alpine:latest


WORKDIR /app

COPY --from=builder /app/publish .

EXPOSE 10001

ENTRYPOINT ["./main"]