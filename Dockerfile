FROM golang:latest
RUN mkdir /server
WORKDIR /server
COPY ./ ./
RUN go env -w GO111MODULE=on
RUN go mod download
RUN go build ./cmd/main/main.go
CMD ["./main"]