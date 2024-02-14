FROM golang:alpine as base-builder

LABEL maintainer='@ctrose17'

WORKDIR /app

EXPOSE 8081
EXPOSE 80

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build

FROM golang:alpine

RUN apk add --no-cache tzdata

WORKDIR /root/

COPY --from=base-builder /app/data/FNameA.csv /root/
COPY --from=base-builder /app/data/FNameB.csv /root/
COPY --from=base-builder /app/data/FNameH.csv /root/
COPY --from=base-builder /app/data/FNameN.csv /root/
COPY --from=base-builder /app/data/FNameW.csv /root/
COPY --from=base-builder /app/data/LNameA.csv /root/
COPY --from=base-builder /app/data/LNameB.csv /root/
COPY --from=base-builder /app/data/LNameH.csv /root/
COPY --from=base-builder /app/data/LNameN.csv /root/
COPY --from=base-builder /app/data/LNameW.csv /root/
COPY --from=0 /app/SimNBA .

ENV PORT 8081
EXPOSE 8081

CMD ["./SimNBA"]