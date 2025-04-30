FROM golang:1.24

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY src/ src/
COPY config config/
RUN go build -o app

CMD ["./app"]
