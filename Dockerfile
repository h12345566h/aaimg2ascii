FROM scratch
WORKDIR $GOPATH/src/github.com/h12345566h/aaimg2ascii
COPY docker  $GOPATH/src/github.com/h12345566h/aaimg2ascii
EXPOSE 8000
CMD ["./aaimg2ascii"]