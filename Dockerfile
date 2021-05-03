FROM golang:latest

LABEL maintainer='@ctrose17 <>'

WORKDIR /app

EXPOSE 8081
EXPOSE 80

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 8081

RUN go build

CMD ["./SimNBA"]