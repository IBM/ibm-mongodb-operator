#!/bin/bash
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
set -e
ICR_NAMESPACE=${ICR_NAMESPACE:-cpopen}
ICR_REPOSITORY=${ICR_REPOSITORY:-ibm-mongodb-operator-app}

[[ "X$ICR_USERNAME" == "X" ]] && read -rp "Enter username icr.io: " ICR_USERNAME
[[ "X$ICR_PASSWORD" == "X" ]] && read -rsp "Enter password icr.io: " ICR_PASSWORD && echo
[[ "X$RELEASE" == "X" ]] && read -rp "Enter Version/Release of operator: " RELEASE

# Fetch authentication token used to push to icr.io
AUTH_TOKEN=$(curl -sH "Content-Type: application/json" -XPOST https://icr.io/cnr/api/v1/users/login -d '
{
    "user": {
        "username": "'"${ICR_USERNAME}"'",
        "password": "'"${ICR_PASSWORD}"'"
    }
}' | awk -F'"' '{print $4}')


# Delete application release in repository
echo "Push package ${ICR_REPOSITORY} into namespace ${ICR_NAMESPACE}"
curl -H "Content-Type: application/json" \
     -H "Authorization: ${AUTH_TOKEN}" \
     -XDELETE https://icr.io/cnr/api/v1/packages/"${ICR_NAMESPACE}"/"${ICR_REPOSITORY}"/"${RELEASE}"/helm
