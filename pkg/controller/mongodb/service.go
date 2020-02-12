package mongodb

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
    app.kubernetes.io/managed-by: helm
    helm.sh/chart: icp-mongodb-3.4.2
    heritage: Helm
    release: mongodb
  name: mongodb
  namespace: ibm-mongodb-operator
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
