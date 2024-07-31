FROM golang:1.20.6

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod .
COPY main.go .
COPY . .
RUN go get
RUN go build -o bin .

ENTRYPOINT ["/app/bin"]