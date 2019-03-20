#!/bin/bash

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

OER_VERSION="v0-9-4"

mvn clean package -DskipTests=true

TMP_BUILD_DIR="oer-collector-release-$OER_VERSION"
# ensure clean build dir
rm -r "$TMP_BUILD_DIR"
mkdir -p "$TMP_BUILD_DIR"
# copy scripts
cp -r release/* "$TMP_BUILD_DIR"/

# get new jars
cp -f oer-collector/target/oer-collector.jar "$TMP_BUILD_DIR"/
cp -f oer-collector-server/target/oer-collector-server.jar "$TMP_BUILD_DIR"/

# bundle tools into release dir
mkdir -p "$TMP_BUILD_DIR"/tools
cp -r tools "$TMP_BUILD_DIR"/tools

chmod +x "$TMP_BUILD_DIR"/start.sh
chmod +x "$TMP_BUILD_DIR"/start_collector.sh
chmod +x "$TMP_BUILD_DIR"/start_server.sh
chmod +x "$TMP_BUILD_DIR"/stop.sh
chmod +x "$TMP_BUILD_DIR"/stop_collector.sh
chmod +x "$TMP_BUILD_DIR"/stop_server.sh
chmod +x "$TMP_BUILD_DIR"/status.sh
chmod +x "$TMP_BUILD_DIR"/restart.sh

# copy license information
cp -f target/generated-sources/license/THIRD-PARTY.txt "$TMP_BUILD_DIR"/THIRD-PARTY.txt
cp -f target/generated-sources/license/THIRD-PARTY.txt THIRD-PARTY.txt
cp -f LICENSE "$TMP_BUILD_DIR"/LICENSE

RELEASE_FILE_NAME="oer-collector-release-$(date -I)-$OER_VERSION"
# TODO implement versioning automatization
zip -r "$RELEASE_FILE_NAME.zip" "$TMP_BUILD_DIR"

# delete conflicting .tar or .tar.gz file
if [ -f "$RELEASE_FILE_NAME.tar" ]; then
    rm "$RELEASE_FILE_NAME.tar"
fi
if [ -f "$RELEASE_FILE_NAME.tar.gz" ]; then
    rm "$RELEASE_FILE_NAME.tar.gz"
fi

tar cf "$RELEASE_FILE_NAME.tar" "$TMP_BUILD_DIR"/
gzip "$RELEASE_FILE_NAME.tar"
rm -r "$TMP_BUILD_DIR"

echo "Generated $RELEASE_FILE_NAME.tar.gz"

echo "Generating SHA256 Checksums:"

sha256sum "$RELEASE_FILE_NAME.tar.gz"
sha256sum "$RELEASE_FILE_NAME.zip"