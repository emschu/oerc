#!/usr/bin/env bash

cd "$(dirname "$0")"
cd spec

# type "make install" if openapi-generator-cli.jar is missing
java -jar ../openapi-generator-cli.jar generate -i openapi3.yaml -g openapi
mv openapi.json openapi3.json

rm -rf ./../.hypothesis
rm -rf ./.openapi-generator
rm -rf ./.openapi-generator-ignore
rm README.md
