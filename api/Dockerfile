FROM golang:1.9

COPY . /go/src/app
WORKDIR /go/src/app

RUN go-wrapper download
RUN go-wrapper install

EXPOSE 80

CMD ["go-wrapper", "run"]
