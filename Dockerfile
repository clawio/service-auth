FROM golang:1.5
MAINTAINER Hugo González Labrador

ENV CLAWIO_AUTH_DBDSN=/tmp/users.db
ENV CLAWIO_AUTH_DBDRIVER=sqlite3
ENV CLAWIO_AUTH_SIGNMETHOD=HS256
ENV CLAWIO_SHAREDSECRET=secret
ENV CLAWIO_AUTH_PORT=57000

ADD . /go/src/github.com/clawio/service-auth
WORKDIR /go/src/github.com/clawio/service-auth

RUN go get -u github.com/tools/godep
RUN godep restore
RUN go install
RUN go get -u github.com/clawio/service-auth

ENTRYPOINT /go/bin/service-auth

EXPOSE 57000

