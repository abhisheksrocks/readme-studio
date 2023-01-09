FROM golang:1.19.4-alpine3.17

# Some basic features
RUN apk --no-cache add git
RUN apk --no-cache add make
RUN apk --no-cache add gcc
RUN apk --no-cache add libc-dev


# The following set of commands will install necessary linting/analysis tools required by IDE
RUN go install -v github.com/ramya-rao-a/go-outline@v0.0.0-20210608161538-9736a4bde949
RUN go install golang.org/x/tools/gopls@v0.11.0
RUN go install github.com/go-delve/delve/cmd/dlv@v1.20.1
RUN go install honnef.co/go/tools/cmd/staticcheck@2022.1.3

WORKDIR /readme-studio

COPY . .

# The following are the program dependencies
RUN go get github.com/joho/godotenv@v1.4.0

CMD go run .