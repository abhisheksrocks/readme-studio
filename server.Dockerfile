FROM golang:1.19.4-alpine3.17

WORKDIR /app

COPY . .

# The following are the program dependencies
RUN go get github.com/joho/godotenv@v1.4.0

CMD go build .;./readme-studio