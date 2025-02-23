FROM golang:1.22.12

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o backend cmd/app/main.go

EXPOSE 8081

CMD ["./backend"]