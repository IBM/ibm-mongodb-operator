#!/usr/bin/env bash

echo "Is the deploy/operator.yaml updated with the latest tags? [Enter once it is]"

EXPORTER_MONGODB_IMAGE=$(cat deploy/operator.yaml | grep mongodb-exporter | cut -d ' ' -f16)
INSTALL_MONGODB_IMAGE=$(cat deploy/operator.yaml | grep mongodb-install | cut -d ' ' -f16)
MONGODB_IMAGE=$(cat deploy/operator.yaml | grep mongodb: | cut -d ' ' -f16)
