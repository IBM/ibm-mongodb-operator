//
// Copyright 2021 IBM Corporation
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

const godIssuerYaml = `
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: god-issuer
  labels:
    app.kubernetes.io/instance: mongodbs.operator.ibm.com
    app.kubernetes.io/managed-by: mongodbs.operator.ibm.com
    app.kubernetes.io/name: mongodbs.operator.ibm.com
spec:
  selfSigned: {}
`

const rootCertYaml = `
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: mongodb-root-ca-cert
  labels:
    app.kubernetes.io/instance: mongodbs.operator.ibm.com
    app.kubernetes.io/managed-by: mongodbs.operator.ibm.com
    app.kubernetes.io/name: mongodbs.operator.ibm.com
spec:
  secretName: mongodb-root-ca-cert
  duration: 17520h
  isCA: true
  issuerRef:
    name: god-issuer
    kind: Issuer
  commonName: "mongodb"
  dnsNames:
  - mongodb.root
`

const rootIssuerYaml = `
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: mongodb-root-ca-issuer
  labels:
    app.kubernetes.io/instance: mongodbs.operator.ibm.com
    app.kubernetes.io/managed-by: mongodbs.operator.ibm.com
    app.kubernetes.io/name: mongodbs.operator.ibm.com
spec:
  ca:
    secretName: mongodb-root-ca-cert
`

const clientCertYaml = `
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: icp-mongodb-client-cert
  labels:
    app.kubernetes.io/instance: mongodbs.operator.ibm.com
    app.kubernetes.io/managed-by: mongodbs.operator.ibm.com
    app.kubernetes.io/name: mongodbs.operator.ibm.com
spec:
  secretName: icp-mongodb-client-cert
  duration: 17520h
  isCA: false
  issuerRef:
    name: mongodb-root-ca-issuer
    kind: Issuer
  commonName: "mongodb-service"
  dnsNames:
  - mongodb
`
