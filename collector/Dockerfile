FROM golang:1.9

COPY . /go/src/app
WORKDIR /go/src/app

RUN go-wrapper download
RUN go-wrapper install

ARG url
ARG countCollectors
ARG countQueueProducers

ENV url ${url}
ENV countCollectors ${countCollectors}
ENV countQueueProducers ${countQueueProducers}

EXPOSE 80

CMD go-wrapper run ${url} ${countCollectors} ${countQueueProducers}
