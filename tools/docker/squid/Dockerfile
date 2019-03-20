FROM alpine:latest
HEALTHCHECK --interval=3s --retries=20 CMD nc -z localhost 3128

RUN apk update && apk add --no-cache squid

EXPOSE 3128/tcp

COPY entrypoint.sh /usr/bin/entrypoint.sh
RUN chmod +x /usr/bin/entrypoint.sh

ENTRYPOINT ["/usr/bin/entrypoint.sh"]

CMD ["squid"]
