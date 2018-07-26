FROM golang:1.10-alpine AS builder

RUN apk add -U make gcc musl-dev ncurses git

ADD .   /go/src/github.com/gnosthi/quixical
WORKDIR /go/src/github.com/gnosthi/quixical

RUN TERM=vt100 make install

FROM alpine:3.7
COPY --from=0 /go/src/github.com/gnosthi/quixical/quixical /usr/bin/
RUN mkdir -p /usr/lib/quixical-problems/
COPY --from=0 /go/src/github.com/gnosthi/quixical/problems/problems-all.csv /usr/lib/quixical-problems/


RUN chown -Rh 1000:1000 -- /root
ENV HOME /root
USER 1000:1000
ENTRYPOINT [ "/usr/bin/quixical","-csv","/usr/lib/quixical-problems/problems-all.csv" ]