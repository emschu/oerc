FROM golang:1.16-alpine

MAINTAINER emschu <emschu@mailbox.org>

RUN mkdir /app
ENV TZ=Europe/Berlin

ADD config/.oerc.docker.yaml /app/.oerc.yaml
ADD bin/oerc-docker /app/oerc
RUN apk add --no-cache tzdata; chmod +x /app/oerc

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["/app/oerc", "-c", "/app/.oerc.yaml"]