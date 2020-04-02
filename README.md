# ibm-mongodb-operator

This operator is built to support a larger service offering from IBM called IBM Cloud Platform Common Services. It stands up a mongoDB database that shared by other services within IBM Cloud Platform Common Services. 

## Supported platforms

This is supported on amd64, ppc64le, and s390x. 

## Operator versions

| Version | Date | Details |
| ----- | ---- | ----------------- |
| 1.0.0 | March 2020 | Initial Offering of MongoDB operator

## Prerequisites

This operator requires cert-manager, https://github.com/IBM/ibm-cert-manager-operator, be installed in the cluster. The cert-manager operator is also part of the IBM Cloud Platform Common Services. You should be using that offering to install both of these operators. 

## Documentation

For installation and configuration, see the [IBM Knowledge Center](http://ibm.biz/cpcsdocs).

Some notes.
The operator does not support updating the CR in version 1.0.0. To make changes to a deployed MongoDB instance it is best to edit the statefulset directly.

When deploying MongoDB, it is better to use 3 replicas, especially if you are not backing up your data. It is possible for the data to become corrupt and recovering from a 3 replica deployment is much easier. 

### Developer guide

Information about building and testing the operator.
- Dev quick start
- Debugging the operator

