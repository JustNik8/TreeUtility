# docker build -t golang_hw1_tree .
FROM golang:1.22
COPY . .
RUN go test -v