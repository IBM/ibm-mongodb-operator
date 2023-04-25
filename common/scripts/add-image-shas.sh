#!/usr/bin/env bash

#
# Copyright 2021 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

echo "Is the deploy/operator.yaml updated with the latest tags? [Enter once it is]"

INSTALL_MONGODB_IMAGE=$(cat deploy/operator.yaml | grep mongodb-install | cut -d ' ' -f16)
MONGODB_IMAGE=$(cat deploy/operator.yaml | grep mongodb: | cut -d ' ' -f16)

INSTALL_SHA=$(docker pull $INSTALL_MONGODB_IMAGE | grep Digest | cut -d ':' -f3)
MONGODB_SHA=$(docker pull $MONGODB_IMAGE | grep Digest | cut -d ':' -f3)

CSV_VERSION=$(cat version/version.go  | grep "Version =" | cut -d '"' -f2)

gsed -i "s/ibm-mongodb-install@sha.*/ibm-mongodb-install@sha256:$INSTALL_SHA/g" deploy/olm-catalog/ibm-mongodb-operator/$CSV_VERSION/ibm-mongodb-operator.v$CSV_VERSION.clusterserviceversion.yaml
gsed -i "s/ibm-mongodb@sha.*/ibm-mongodb@sha256:$MONGODB_SHA/g" deploy/olm-catalog/ibm-mongodb-operator/$CSV_VERSION/ibm-mongodb-operator.v$CSV_VERSION.clusterserviceversion.yaml
