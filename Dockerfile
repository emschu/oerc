#
# oerc, alias oer-collector
# Copyright (C) 2021 emschu[aet]mailbox.org
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public
# License along with this program.
# If not, see <https://www.gnu.org/licenses/>.

# Build
FROM golang:1.17-alpine as build
RUN mkdir /app
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download && mkdir /app/project
COPY . ./project
RUN cd ./project && go build -o /app/oerc


FROM golang:1.17-alpine
MAINTAINER emschu <emschu@mailbox.org>
RUN mkdir /app && apk add --no-cache tzdata;
WORKDIR /app
ENV TZ=Europe/Berlin

COPY --from=build /app/oerc /app/oerc
RUN chmod +x /app/oerc
EXPOSE 8080
ADD config/.oerc.docker.yaml /app/.oerc.yaml
ENTRYPOINT ["/app/oerc", "-c", "/app/.oerc.yaml"]