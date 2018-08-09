FROM golang:alpine
MAINTAINER Transceptor Technology
RUN apk add git
WORKDIR /src
RUN git clone https://github.com/SiriDB/siridb-prompt ./
RUN go get -d
RUN go build -o siridb-prompt

FROM alpine:latest
COPY --from=0 /src/siridb-prompt /usr/local/bin/siridb-prompt
ENTRYPOINT ["/usr/local/bin/siridb-prompt"]
