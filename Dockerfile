FROM golang:alpine as base-builder

LABEL maintainer='@ctrose17'

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build

FROM golang:alpine

RUN apk add --no-cache tzdata

WORKDIR /root/

COPY --from=0 /app/SimNBA .

COPY --from=base-builder /app/data /app/data


ENV PORT 8081
ENV ROOT=/app
ENV GOPATH /go
EXPOSE 8081

CMD ["./SimNBA"]