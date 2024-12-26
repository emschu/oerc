#!/usr/bin/env bash
#
# oerc, alias oer-collector
# Copyright (C) 2021-2024 emschu[aet]mailbox.org
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

cd "$(dirname "$0")"
cd spec

# type "make setup" if openapi-generator-cli.jar is missing
java -jar ../openapi-generator-cli.jar generate -i openapi3.yaml -g openapi
mv openapi.json openapi3.json

rm -rf ./../.hypothesis
rm -rf ./.openapi-generator
rm -rf ./.openapi-generator-ignore
rm README.md
