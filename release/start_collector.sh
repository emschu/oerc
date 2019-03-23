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

touch collector.log

echo "Starting OER Collector component..."

# use user defined java
if [ -z "$JAVA_HOME" ]; then
    JAVA_EXEC="java"
else
    JAVA_EXEC="$JAVA_HOME/java"
    echo "using java: JAVA_EXEC"
fi

# export app properties
SPRING_APPLICATION_JSON=$(cat config.json)
echo $SPRING_APPLICATION_JSON
export SPRING_APPLICATION_JSON

# build java cmd
CMD="$JAVA_EXEC -Xmx164M -Xms164M -jar oer-collector.jar"
echo "$CMD"

# check for foreground mode
if [[ $* == *--foreground* ]]; then
    echo "Running in foreground"
    command $CMD
    echo "Collector finished"
    exit 0
fi

# use nohup if possible
if ! [ -x "$(command -v nohup)" ]; then
    $($CMD >> collector.log 2>&1) &
else
    nohup $CMD >> collector.log 2>&1 &
fi
echo "Collector started in background process: " $(jps -l | grep "oer-collector\.jar" | awk '{print $1}')
exit 0