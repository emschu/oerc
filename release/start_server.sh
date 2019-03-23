#!/usr/bin/env bash

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

touch server.log

echo "Starting OER Collector REST API server..."

# use user defined java
if [ -z "$JAVA_HOME" ]; then
    JAVA_EXEC="java"
else
    JAVA_EXEC="$JAVA_HOME/bin/java"
    echo "using java: JAVA_EXEC"
fi

# export app properties
SPRING_APPLICATION_JSON=$(cat config.json)
export SPRING_APPLICATION_JSON

# build java cmd
CMD="$JAVA_EXEC -Xmx64M -Xmx64M -jar oer-collector-server.jar $JAVA_ARGS"
echo "$CMD"

# check if already running
if [[ $(jps -l | grep "oer-collector-server.jar" | awk '{print $1}') ]]; then
    echo "OER Server is already running!"
    exit 1
fi

# check for foreground mode
if [[ $* == *--foreground* ]]; then
    echo "Running in foreground"
    command $CMD
    echo "OER Server finished"
    exit 0
fi

# use nohup if possible
if ! [ -x "$(command -v nohup)" ]; then
    $($CMD >> server.log 2>&1) &
else
    nohup $CMD >> server.log 2>&1 &
fi
echo "Server started in background process: " $(jps -l | grep "oer-collector-server.jar" | awk '{print $1}')
exit 0