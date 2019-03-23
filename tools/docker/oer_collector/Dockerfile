FROM openjdk:8-jdk-alpine

COPY oer-collector.jar /oer-collector.jar
COPY start_collector.sh /start_collector.sh
RUN apk add bash; chmod +x /start_collector.sh
COPY stop_collector.sh /stop_collector.sh
COPY config.json /config.json

CMD ["bash", "start_collector.sh", "--foreground"]