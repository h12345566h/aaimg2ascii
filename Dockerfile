FROM scratch
WORKDIR $GOPATH/src/github.com/h12345566h/aaimg2ascii
COPY . $GOPATH/src/github.com/h12345566h/aaimg2ascii
EXPOSE 80
CMD ["./aaimg2ascii"]