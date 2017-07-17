FROM golang:1.8
RUN echo 'set -o vi' >> /root/.bashrc

# Fetch dependencies
RUN go get github.com/tools/godep

ENV USER topher
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET ittwiP92o0oi6P4i
ENV DSN postgres://postgres@db:5432/forty-thieves?sslmode=disable

ADD . /go/src/github.com/topher200/forty-thieves
WORKDIR /go/src/github.com/topher200/forty-thieves

RUN godep go build

EXPOSE 8888
CMD ./forty-thieves
