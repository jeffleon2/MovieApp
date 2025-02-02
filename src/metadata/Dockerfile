FROM golang:1.19
WORKDIR /go/src/
COPY main .
COPY configs/. .
EXPOSE 8083
EXPOSE 8500
RUN chmod a+x main
CMD [ "./main" ]