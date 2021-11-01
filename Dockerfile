FROM golang:1.17-alpine as build
WORKDIR /app
COPY . .
RUN go build -o effrit main.go

WORKDIR /app/
ENTRYPOINT [ "/app/effrit" ]
