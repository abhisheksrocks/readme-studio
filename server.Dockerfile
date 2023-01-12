FROM golang:1.19.5-alpine3.17

WORKDIR /app

COPY . .

# The following are the program dependencies
RUN go get github.com/joho/godotenv@v1.4.0
RUN go get github.com/stretchr/testify@v1.8.1

CMD go build .;./readme-studio