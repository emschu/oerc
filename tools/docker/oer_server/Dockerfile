FROM openjdk:8-jdk-alpine

COPY oer-collector-server.jar /oer-collector-server.jar
COPY start_server.sh /start_server.sh
RUN apk add bash; chmod +x start_server.sh
COPY stop_server.sh /stop_server.sh
COPY config.json /config.json

EXPOSE 8081

CMD ["bash", "start_server.sh", "--foreground"]