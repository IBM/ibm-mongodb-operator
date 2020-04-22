# ibm-mongodb-operator

This operator is built to support a larger service offering from IBM called IBM Cloud Platform Common Services. It stands up a mongoDB database that shared by other services within IBM Cloud Platform Common Services.

## Supported platforms

This is supported on amd64, ppc64le, and s390x.

## Operator versions

| Version | Date | Details |
| ----- | ---- | ----------------- |
| 1.1.0 | TDB (In Development) | Allow user to bring their own admin secret
| 1.0.0 | March 2020 | Initial Offering of MongoDB operator

## Prerequisites

This operator requires cert-manager, https://github.com/IBM/ibm-cert-manager-operator, be installed in the cluster. The cert-manager operator is also part of the IBM Cloud Platform Common Services. You should be using that offering to install both of these operators.

## Documentation

For installation and configuration, see the [IBM Knowledge Center](http://ibm.biz/cpcsdocs).

#### Highlighted Features

**_Admin Secret_**

Starting with version 1.1.0 you can now supply your own `icp-mongodb-admin` secret. The secret must have a `user` field and a `password` field and be in the same namespace that mongoDB is going to be created in. If the user chooses not to supply a secret, a random user and password will be created and used. The `icp-mongodb-admin` secret will persist after uninstalling/removing the MongoDB custom resource so that uninstall and re-install is possible using the same Persistent Volumes.

Example Yaml for creating your own admin secret before installation. The user and password are base64 encrypted.
```
apiVersion: v1
kind: Secret
metadata:
  name: icp-mongodb-admin
  namespace: ibm-common-services
type: Opaque
data:
  password: SFV6a2NYMkdKa2tBZA==
  user: dGpOcDR5Unc=
```

#### Notes
The operator does not support updating the CR in version 1.0.0. To make changes to a deployed MongoDB instance it is best to edit the statefulset directly.

When deploying MongoDB, it is better to use 3 replicas, especially if you are not backing up your data. It is possible for the data to become corrupt and recovering from a 3 replica deployment is much easier.

### What's New

In version 1.1.0
- Allow user to bring their own admin secret
- The CSV defines dependencies it has to run


### Developer guide

Information about building and testing the operator.
- Dev quick start (TO-DO)
- Debugging the operator (TO-DO)
