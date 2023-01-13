FROM golang:1.19.5-alpine3.17

WORKDIR /app

COPY . .

CMD go build .;./readme-studio