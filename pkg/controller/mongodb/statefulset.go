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

const statefulset = `
---
# Source: icp-mongodb/templates/mongodb-statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    app: icp-mongodb
    release: mongodb
    app.kubernetes.io/instance: mongodbs.operator.ibm.com
    app.kubernetes.io/managed-by: mongodbs.operator.ibm.com
    app.kubernetes.io/name: mongodbs.operator.ibm.com
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
        app.kubernetes.io/instance: common-mongodb
        app: icp-mongodb
        release: mongodb
      annotations:
        productName: "IBM Cloud Platform Common Services"
        productID: "068a62892a1e4db39641342e592daa25"
        productVersion: "3.4.0"
        productMetric: "FREE"
        prometheus.io/scrape: "true"
        prometheus.io/port: "9216"
        prometheus.io/path: "/metrics"
        clusterhealth.ibm.com/dependencies: ibm-common-services.cert-manager
    spec:
      serviceAccountName: ibm-mongodb-operand
      terminationGracePeriodSeconds: 30
      hostNetwork: false
      hostPID: false
      hostIPC: false
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/arch
                operator: In
                values:
                  - amd64
                  - ppc64le
                  - s390x
      initContainers:
        - name: install
          image: "{{ .InitImage }}"
          command:
            - /install/install.sh
          args:
            - --work-dir=/work-dir
            - --config-dir=/data/configdb
          imagePullPolicy: "IfNotPresent"
          resources:
            limits:
              cpu: {{ .CpuLimit }}
              memory: {{ .MemoryLimit }}
            requests:
              cpu: {{ .CpuRequest }}
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
          imagePullPolicy: "IfNotPresent"
          resources:
            limits:
              cpu: {{ .CpuLimit }}
              memory: {{ .MemoryLimit }}
            requests:
              cpu: {{ .CpuRequest }}
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
          imagePullPolicy: "IfNotPresent"
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
          ports:
            - name: peer
              containerPort: 27017
          resources:
            limits:
              cpu: {{ .CpuLimit }}
              memory: {{ .MemoryLimit }}
            requests:
              cpu: {{ .CpuRequest }}
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
                - mongo
                - --ssl
                - --sslCAFile=/data/configdb/tls.crt
                - --sslPEMKeyFile=/work-dir/mongo.pem
                - --eval
                - "db.adminCommand('ping')"
            initialDelaySeconds: 30
            timeoutSeconds: 5
            failureThreshold: 3
            periodSeconds: 10
            successThreshold: 1
          readinessProbe:
            exec:
              command:
                - mongo
                - --ssl
                - --sslCAFile=/data/configdb/tls.crt
                - --sslPEMKeyFile=/work-dir/mongo.pem
                - --eval
                - "db.adminCommand('ping')"
            initialDelaySeconds: 5
            timeoutSeconds: 1
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

        - name: metrics
          image: "{{ .MetricsImage }}"
          imagePullPolicy: "IfNotPresent"
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
          command:
            - sh
            - -ec
            - >-
              /bin/mongodb_exporter
              --mongodb.uri mongodb://$METRICS_USER:$METRICS_PASSWORD@localhost:27017
              --mongodb.tls
              --mongodb.tls-ca=/data/configdb/tls.crt
              --mongodb.tls-cert=/work-dir/mongo.pem
              --mongodb.socket-timeout=3s
              --mongodb.sync-timeout=1m
              --web.telemetry-path=/metrics
              --web.listen-address=:9216
          volumeMounts:
            - name: configdir
              mountPath: /data/configdb
            - name: mongodbdir
              subPath: workdir
              mountPath: /work-dir
            - name: tmp-metrics
              mountPath: /tmp
          env:
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
          ports:
            - name: metrics
              containerPort: 9216
          resources:
            limits:
              cpu: 500m
              memory: 256Mi
          readinessProbe:
            exec:
              command:
                - sh
                - -ec
                - >-
                  /bin/mongodb_exporter
                  --mongodb.uri mongodb://$METRICS_USER:$METRICS_PASSWORD@localhost:27017
                  --mongodb.tls
                  --mongodb.tls-ca=/data/configdb/tls.crt
                  --mongodb.tls-cert=/work-dir/mongo.pem
                  --test
            initialDelaySeconds: 30
            periodSeconds: 10
          livenessProbe:
            exec:
              command:
                - sh
                - -ec
                - >-
                  /bin/mongodb_exporter
                  --mongodb.uri mongodb://$METRICS_USER:$METRICS_PASSWORD@localhost:27017
                  --mongodb.tls
                  --mongodb.tls-ca=/data/configdb/tls.crt
                  --mongodb.tls-cert=/work-dir/mongo.pem
                  --test
            initialDelaySeconds: 30
            periodSeconds: 10
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
            storage: 20Gi
        storageClassName: {{ .StorageClass }}
`
