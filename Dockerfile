FROM golang:alpine as base-builder

LABEL maintainer='@ctrose17 <>'

WORKDIR /app

EXPOSE 8081
EXPOSE 80

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /app/SimNBA .

ENV PORT 8081
EXPOSE 8081

CMD ["./SimNBA"]