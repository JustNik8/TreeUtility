# docker build -t golang_hw1_tree .
FROM golang:1.22-alpine
COPY . .
CMD ["go","test", "-v"]