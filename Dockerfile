FROM golang:1.17-alpine as build
WORKDIR /app
COPY . .
RUN go build -o effrit main.go

FROM alpine
RUN apk add -u ca-certificates
COPY --from=build /app/effrit /app/

WORKDIR /app/
ENTRYPOINT [ "/app/effrit" ]
