FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/topher200/forty-thieves

ENV USER topher
ENV HTTP_ADDR 8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET ittwiP92o0oi6P4i

# Replace this with actual PostgreSQL DSN.
ENV DSN postgres://topher@localhost:5432/forty-thieves?sslmode=disable

WORKDIR /go/src/github.com/topher200/forty-thieves

RUN godep go build

EXPOSE 8888
CMD ./forty-thieves
