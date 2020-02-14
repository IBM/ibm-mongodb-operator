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

package mongodb

const mongodbSCC = `
allowHostDirVolumePlugin: true
allowHostIPC: false
allowHostNetwork: true
allowHostPID: true
allowHostPorts: true
allowPrivilegedContainer: false
allowedCapabilities: null
allowedFlexVolumes: null
apiVersion: security.openshift.io/v1
defaultAddCapabilities: null
fsGroup:
  type: RunAsAny
users:
  - system:serviceaccount:ibm-mongodb-operator:default
kind: SecurityContextConstraints
metadata:
  name: ibm-mongodb-scc
priority: 1
readOnlyRootFilesystem: false
requiredDropCapabilities:
  - MKNOD
runAsUser:
  type: RunAsAny
seLinuxContext:
  type: MustRunAs
seccompProfiles:
- '*'
supplementalGroups:
  type: RunAsAny
groups: []
volumes:
  - configMap
  - downwardAPI
  - emptyDir
  - hostPath
  - nfs
  - persistentVolumeClaim
  - projected
  - secret
`
