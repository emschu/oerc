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

if [ "$(jps -l | grep -c 'oer-collector')" -gt 0 ]; then
    echo "Killing $(jps -l | grep -c 'oer-collector') collector processes"
    jps -vl | grep 'oer-collector'
    jps -l | grep "oer-collector" | awk '{print $1}' | xargs kill -9
else
    echo "No running OER Collector processes found"
fi
