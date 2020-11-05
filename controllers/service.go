//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package controllers

const service = `
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/name: icp-mongodb
    app.kubernetes.io/instance: icp-mongodb
    app.kubernetes.io/version: 4.0.12-build.3
    app.kubernetes.io/component: database
    app.kubernetes.io/part-of: common-services-cloud-pak
    app.kubernetes.io/managed-by: operator
    release: mongodb
  name: mongodb
spec:
  serviceAccountName: ibm-mongodb-operator
  type: ClusterIP
  ports:
  - port: 27017
    protocol: TCP
    targetPort: 27017
  selector:
    app: icp-mongodb
    release: mongodb
`
