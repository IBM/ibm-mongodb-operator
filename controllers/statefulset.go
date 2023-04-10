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

const statefulset = `
---
# Source: icp-mongodb/templates/mongodb-statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    {{- range $key, $value := .StsLabels }}
      {{ $key}}: {{ $value}}
    {{- end }}
  name: icp-mongodb
spec:
  selector:
    matchLabels:
      app: icp-mongodb
      release: mongodb
  serviceName: icp-mongodb
  replicas: {{ .Replicas }}
  template:
    metadata:
      labels:
        {{- range $key, $value := .PodLabels }}
          {{ $key}}: {{ $value}}
        {{- end }}
      annotations:
        productName: "IBM Cloud Platform Common Services"
        productID: "068a62892a1e4db39641342e592daa25"
        productMetric: "FREE"
        prometheus.io/scrape: "true"
        prometheus.io/port: "9216"
        prometheus.io/path: "/metrics"
        clusterhealth.ibm.com/dependencies: {{ .NamespaceName }}.cert-manager
    spec:
      serviceAccountName: ibm-mongodb-operand
      {{ if eq .UserID 1000 }}
      securityContext:
        runAsUser: {{ .UserID }}
        fsGroup: 0
      {{ end }}
      terminationGracePeriodSeconds: 30
      hostNetwork: false
      hostPID: false
      hostIPC: false
      topologySpreadConstraints:
      - maxSkew: 1
        topologyKey: topology.kubernetes.io/zone
        whenUnsatisfiable: ScheduleAnyway
        labelSelector:
          matchLabels:
            key: app
            values: icp-mongodb
      - maxSkew: 1
        topologyKey: topology.kubernetes.io/region
        whenUnsatisfiable: ScheduleAnyway
        labelSelector:
          matchLabels:
            key: app
            values: icp-mongodb
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 50
            podAffinityTerm:
              topologyKey: kubernetes.io/hostname
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - icp-mongodb
      initContainers:
        - name: install
          image: "{{ .InitImage }}"
          command:
            - /install/install.sh
          args:
            - --work-dir=/work-dir
            - --config-dir=/data/configdb
          imagePullPolicy: "Always"
          resources:
            limits:
              cpu: {{ .CPULimit }}
              memory: {{ .MemoryLimit }}
            requests:
              cpu: {{ .CPURequest }}
              memory: {{ .MemoryRequest }}
          volumeMounts:
            - name: mongodbdir
              subPath: workdir
              mountPath: /work-dir
            - name: configdir
              mountPath: /data/configdb
            - name: config
              mountPath: /configdb-readonly
            - name: install
              mountPath: /install
            - name: keydir
              mountPath: /keydir-readonly
            - name: ca
              mountPath: /ca-readonly
            - name: mongodbdir
              subPath: datadir
              mountPath: /data/db
            - name: tmp-mongodb
              mountPath: /tmp
        - name: bootstrap
          image: "{{ .BootstrapImage }}"
          command:
            - /work-dir/peer-finder
          args:
            - -on-start=/init/on-start.sh
            - "-service=icp-mongodb"
          imagePullPolicy: "Always"
          resources:
            limits:
              cpu: {{ .CPULimit }}
              memory: {{ .MemoryLimit }}
            requests:
              cpu: {{ .CPURequest }}
              memory: {{ .MemoryRequest }}
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
          env:
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  apiVersion: v1
                  fieldPath: metadata.namespace
            - name: REPLICA_SET
              value: rs0
            - name: AUTH
              value: "true"
            - name: ADMIN_USER
              valueFrom:
                secretKeyRef:
                  name: "icp-mongodb-admin"
                  key: user
            - name: ADMIN_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: "icp-mongodb-admin"
                  key: password
            - name: METRICS
              value: "true"
            - name: METRICS_USER
              valueFrom:
                secretKeyRef:
                  name: "icp-mongodb-metrics"
                  key: user
            - name: METRICS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: "icp-mongodb-metrics"
                  key: password
            - name: NETWORK_IP_VERSION
              value: ipv4
          volumeMounts:
            - name: mongodbdir
              subPath: workdir
              mountPath: /work-dir
            - name: configdir
              mountPath: /data/configdb
            - name: init
              mountPath: /init
            - name: mongodbdir
              subPath: datadir
              mountPath: /data/db
            - name: tmp-mongodb
              mountPath: /tmp
      containers:
        - name: icp-mongodb
          image: "{{ .BootstrapImage }}"
          imagePullPolicy: "Always"
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
          ports:
            - name: peer
              containerPort: 27017
          resources:
            limits:
              cpu: {{ .CPULimit }}
              memory: {{ .MemoryLimit }}
            requests:
              cpu: {{ .CPURequest }}
              memory: {{ .MemoryRequest }}
          command:
            - mongod
            - --config=/data/configdb/mongod.conf
          env:
            - name: AUTH
              value: "true"
            - name: ADMIN_USER
              valueFrom:
                secretKeyRef:
                  name: "icp-mongodb-admin"
                  key: user
            - name: ADMIN_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: "icp-mongodb-admin"
                  key: password
          livenessProbe:
            exec:
              command:
                - mongosh
                - --tls
                - --tlsCAFile=/data/configdb/tls.crt
                - --tlsCertificateKeyFile=/work-dir/mongo.pem
                - --eval
                - "db.adminCommand('ping')"
            initialDelaySeconds: 30
            timeoutSeconds: 10
            failureThreshold: 5
            periodSeconds: 30
            successThreshold: 1
          readinessProbe:
            exec:
              command:
                - mongosh
                - --tls
                - --tlsCAFile=/data/configdb/tls.crt
                - --tlsCertificateKeyFile=/work-dir/mongo.pem
                - --eval
                - "db.adminCommand('ping')"
            initialDelaySeconds: 5
            timeoutSeconds: 5
            failureThreshold: 3
            periodSeconds: 10
            successThreshold: 1
          volumeMounts:
            - name: mongodbdir
              subPath: datadir
              mountPath: /data/db
            - name: configdir
              mountPath: /data/configdb
            - name: mongodbdir
              subPath: workdir
              mountPath: /work-dir
            - name: tmp-mongodb
              mountPath: /tmp
      tolerations:
        - effect: NoSchedule
          key: dedicated
          operator: Exists
        - key: CriticalAddonsOnly
          operator: Exists
        - effect: NoExecute
          key: node.kubernetes.io/not-ready
          operator: Exists
        - effect: NoExecute
          key: node.kubernetes.io/unreachable
          operator: Exists
      volumes:
        - name: config
          configMap:
            name: icp-mongodb
        - name: init
          configMap:
            defaultMode: 0755
            name: icp-mongodb-init
        - name: install
          configMap:
            defaultMode: 0755
            name: icp-mongodb-install
        - name: ca
          secret:
            defaultMode: 0755
            secretName: mongodb-root-ca-cert
        - name: keydir
          secret:
            defaultMode: 0755
            secretName: icp-mongodb-keyfile
        - name: configdir
          emptyDir: {}
        - name: tmp-mongodb
          emptyDir: {}
        - name: tmp-metrics
          emptyDir: {}
  volumeClaimTemplates:
    - metadata:
        name: mongodbdir
      spec:
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: {{ .PVCSize }}
        storageClassName: {{ .StorageClass }}
`
