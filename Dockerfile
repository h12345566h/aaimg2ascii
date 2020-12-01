FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/github.com/h12345566h/aaimg2ascii
COPY . $GOPATH/src/github.com/h12345566h/aaimg2ascii
RUN go build .
EXPOSE 8000

ENTRYPOINT ["./aaimg2ascii"]