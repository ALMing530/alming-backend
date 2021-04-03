FROM golang:1.16.3-alpine3.13
ENV GOPROXY="https://goproxy.cn, direct"
ENV GO111MODULE=on
ENV GOPATH=/home/alming/GOPATH
WORKDIR /home/alming/go
COPY ./ ./alming_backend
WORKDIR /home/alming/go/alming_backend
RUN go mod download
RUN go build ./src/main.go
EXPOSE 53000
CMD [ "./main" ]
