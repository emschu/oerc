#!/bin/#!/usr/bin/env bash

###
# #%L
# oer-collector-project
# %%
# Copyright (C) 2019 emschu[aet]mailbox.org
# %%
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
# 
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
# #L%
###

# script parameters:

SOLR_FOLDER=~/software/solr-7.7.0

if [ ! -d  "$SOLR_FOLDER" ]; then
    echo "solr folder $SOLR_FOLDER does not exist. Please edit script."
    exit 1
fi

# changing this var is not enough to change mysql db driver version this script uses!
DRIVER_NAME=mariadb-java-client-2.4.0.jar

if [ ! -f "$SOLR_FOLDER/lib/$DRIVER_NAME" ]; then
  wget http://central.maven.org/maven2/org/mariadb/jdbc/mariadb-java-client/2.4.0/"$DRIVER_NAME"
  mkdir -p "$SOLR_FOLDER/lib/"
  mv mariadb-java-client-*.jar "$SOLR_FOLDER/lib/"
fi

rm -r "$SOLR_FOLDER"/server/solr/oer-server

"$SOLR_FOLDER"/bin/solr stop
"$SOLR_FOLDER"/bin/solr start
"$SOLR_FOLDER"/bin/solr delete -c oer-server
"$SOLR_FOLDER"/bin/solr create -c oer-server -shards 2 -replicationFactor 2
"$SOLR_FOLDER"/bin/solr stop

echo "start configuration"

sed -i '/\/queryResponseWriter/a\
   <requestHandler name="/dataimport" class="org.apache.solr.handler.dataimport.DataImportHandler"> \
    <lst name="defaults"> \
        <str name="config">data-config.xml</str> \
    </lst> \
  </requestHandler>' "$SOLR_FOLDER"/server/solr/oer-server/conf/solrconfig.xml

sed -i '/luceneMatchVersion/ a\
<lib dir="${solr.install.dir:../../../..}/lib/" regex="mariadb-java-client-2.4.0.jar" />\
<lib dir="${solr.install.dir:../../../..}/dist/" regex="solr-dataimporthandler-.*\.jar" />\
<lib dir="${solr.install.dir:../../../..}/contrib/extraction/lib" regex=".*\.jar" />' "$SOLR_FOLDER"/server/solr/oer-server/conf/solrconfig.xml

cat data-config.xml > "$SOLR_FOLDER"/server/solr/oer-server/conf/data-config.xml

if [ -f "$SOLR_FOLDER/server/solr/oer-server/conf/dataimport.properties" ]; then
    chmod 777 "$SOLR_FOLDER"/server/solr/oer-server/conf/dataimport.properties
fi

echo "Starting Apache Solr..."
"$SOLR_FOLDER"/bin/solr start

echo "Add fields via solr api"
ENDPOINT=http://localhost:8983/solr/oer-server/schema

curl -X POST -H 'Content-type:application/json' --data-binary '{"add-field": {"name":"title", "type":"text_general", "multiValued":false, "stored":true}}' "$ENDPOINT"
curl -X POST -H 'Content-type:application/json' --data-binary '{"add-field": {"name":"description", "type":"text_general", "multiValued":false, "stored":true}}' "$ENDPOINT"
curl -X POST -H 'Content-type:application/json' --data-binary '{"add-field": {"name":"start_date_time", "type":"pdate", "multiValued":false, "stored":true}}' "$ENDPOINT"
curl -X POST -H 'Content-type:application/json' --data-binary '{"add-field": {"name":"end_date_time", "type":"pdate", "multiValued":false, "stored":true}}' "$ENDPOINT"
curl -X POST -H 'Content-type:application/json' --data-binary '{"add-field": {"name":"duration", "type":"pint", "multiValued":false, "stored":true}}' "$ENDPOINT"
