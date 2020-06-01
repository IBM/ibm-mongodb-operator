#!/usr/bin/env bash

#
# Copyright 2020 IBM Corporation
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

# This script needs to inputs
# The CSV version that is currently in dev

CURRENT_DEV_CSV=$1
let NEW_DEV_CSV_Z=$(echo $CURRENT_DEV_CSV | cut -d '.' -f3)+1
NEW_DEV_CSV=$(echo $CURRENT_DEV_CSV | gsed "s/\.[0-9][0-9]*$/\.$NEW_DEV_CSV_Z/")

CSV_PATH=deploy/olm-catalog/ibm-mongodb-operator/
#echo $NEW_DEV_CSV
# Make new z level release directory
mkdir $CSV_PATH/$NEW_DEV_CSV
echo "Made new directory"
read
# Copy Current CSV directory to new one
cp $CSV_PATH/$CURRENT_DEV_CSV/* $CSV_PATH/$NEW_DEV_CSV/
echo "Copied current csv to new directory"
read

# Change to new CSV Version
mv $CSV_PATH/$NEW_DEV_CSV/ibm-mongodb-operator.v$CURRENT_DEV_CSV.clusterserviceversion.yaml $CSV_PATH/$NEW_DEV_CSV/ibm-mongodb-operator.v$NEW_DEV_CSV.clusterserviceversion.yaml
echo "Changed file name csv in new directory"
read

# Update New CSV
# replace old CSV value with new one
gsed -i "s/$CURRENT_DEV_CSV/$NEW_DEV_CSV/g" $CSV_PATH/$NEW_DEV_CSV/ibm-mongodb-operator.v$NEW_DEV_CSV.clusterserviceversion.yaml
TIME_STAMP=$(date '+%Y-%m-%dT%H:%M:%S'Z)
gsed -i "s/2[0-9]*-[0-9]*-[0-9]*T[0-9]*:[0-9]*:[0-9]*Z/$TIME_STAMP/g" $CSV_PATH/$NEW_DEV_CSV/ibm-mongodb-operator.v$NEW_DEV_CSV.clusterserviceversion.yaml
echo "Updated New file with new CSV version"
read

#Update old CSV
# Get SHA values for all images and replace tags
MONGODB_OPERATOR_SHA=$(docker pull quay.io/opencloudio/ibm-mongodb-operator:$CURRENT_DEV_CSV | grep Digest | cut -d ' ' -f2)
MONGODB_INSTALL_TAG=$(cat $CSV_PATH/$CURRENT_DEV_CSV/ibm-mongodb-operator.v$CURRENT_DEV_CSV.clusterserviceversion.yaml | grep value: | grep mongodb-install | cut -d ':' -f3)
MONGODB_EXPORTER_TAG=$(cat $CSV_PATH/$CURRENT_DEV_CSV/ibm-mongodb-operator.v$CURRENT_DEV_CSV.clusterserviceversion.yaml | grep value: | grep mongodb-exporter | cut -d ':' -f3)
MONGODB_IMAGE_TAG=$(cat $CSV_PATH/$CURRENT_DEV_CSV/ibm-mongodb-operator.v$CURRENT_DEV_CSV.clusterserviceversion.yaml | grep value: | grep mongodb: | cut -d ':' -f3)

MONGODB_INSTALL_SHA=$(docker pull quay.io/opencloudio/ibm-mongodb-install:$MONGODB_INSTALL_TAG | grep Digest | cut -d ' ' -f2)
MONGODB_EXPORTER_SHA=$(docker pull quay.io/opencloudio/ibm-mongodb-exporter:$MONGODB_EXPORTER_TAG | grep Digest | cut -d ' ' -f2)
MONGODB_IMAGE_SHA=$(docker pull quay.io/opencloudio/ibm-mongodb:$MONGODB_IMAGE_TAG | grep Digest | cut -d ' ' -f2)

gsed -i "s/ibm-mongodb-install:$MONGODB_INSTALL_TAG/ibm-mongodb-install@$MONGODB_INSTALL_SHA/g" $CSV_PATH/$CURRENT_DEV_CSV/ibm-mongodb-operator.v$CURRENT_DEV_CSV.clusterserviceversion.yaml
gsed -i "s/ibm-mongodb-exporter:$MONGODB_EXPORTER_TAG/ibm-mongodb-exporter@$MONGODB_EXPORTER_SHA/g" $CSV_PATH/$CURRENT_DEV_CSV/ibm-mongodb-operator.v$CURRENT_DEV_CSV.clusterserviceversion.yaml
gsed -i "s/ibm-mongodb:$MONGODB_IMAGE_TAG/ibm-mongodb@$MONGODB_IMAGE_SHA/g" $CSV_PATH/$CURRENT_DEV_CSV/ibm-mongodb-operator.v$CURRENT_DEV_CSV.clusterserviceversion.yaml
gsed -i "s/ibm-mongodb-operator:latest/ibm-mongodb-operator@$MONGODB_OPERATOR_SHA/g" $CSV_PATH/$CURRENT_DEV_CSV/ibm-mongodb-operator.v$CURRENT_DEV_CSV.clusterserviceversion.yaml

echo "Updated current CSV with SHAs from quay"
read

# Update Package channels beta and dev
echo "MANUALLY UPDATE the package.yaml with new channel CSVs (Push Enter when YOU HAVE MANUALLY UPDATED the package): "
read

#Update version.go to new dev version
gsed -i "s/$CURRENT_DEV_CSV/$NEW_DEV_CSV/" version/version.go
gsed -i "s/$CURRENT_DEV_CSV/$NEW_DEV_CSV/" Makefile
echo "Updated the version.go with new version (Push Enter when done): "
read

# Push CSV package yaml to quay
common/scripts/push-csv.sh
echo "Pushed CSV to quay "
read
