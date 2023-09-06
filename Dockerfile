FROM golang:1.19-bookworm
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN ENV=prod go build -o ./out/dist .
CMD ./out/dist