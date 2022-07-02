FROM golang:1.18.2

WORKDIR /server

LABEL project="ascii-art-web" devs="zazuna, koyanchoz" version="1.0"

COPY . .

RUN go build -o server -v ./server/server.go

ENTRYPOINT [ "./server/server" ] 