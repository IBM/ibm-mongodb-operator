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

import (
	"context"

	operatorv1alpha1 "github.com/IBM/ibm-mongodb-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type statusRetrievalFunc func(context.Context, client.Client, []string, string) []operatorv1alpha1.ManagedResourceStatus

const (
	UnknownAPIVersion     string = "Unknown"
	ResourceReadyState    string = "Ready"
	ResourceNotReadyState string = "NotReady"
)

func getServiceStatus(ctx context.Context, k8sClient client.Client, namespacedName types.NamespacedName) (status operatorv1alpha1.ManagedResourceStatus) {
	reqLogger := logf.FromContext(ctx).WithName("getServiceStatus")
	kind := "Service"
	status = operatorv1alpha1.ManagedResourceStatus{
		ObjectName: namespacedName.Name,
		APIVersion: UnknownAPIVersion,
		Namespace:  namespacedName.Namespace,
		Kind:       kind,
		Status:     ResourceNotReadyState,
	}
	service := &corev1.Service{}
	err := k8sClient.Get(ctx, namespacedName, service)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Could not find resource for status update", "kind", kind, "name", namespacedName.Name, "namespace", namespacedName.Namespace)
		} else {
			reqLogger.Error(err, "Error reading resource for status update", "kind", kind, "name", namespacedName.Name, "namespace", namespacedName.Namespace)
		}
		return
	}
	status.APIVersion = service.APIVersion
	status.Status = ResourceReadyState
	return
}

func getAllServiceStatus(ctx context.Context, k8sClient client.Client, names []string, namespace string) (statuses []operatorv1alpha1.ManagedResourceStatus) {
	reqLogger := logf.FromContext(ctx).WithName("getAllServiceStatus").V(3)
	for _, name := range names {
		nsn := types.NamespacedName{Name: name, Namespace: namespace}
		statuses = append(statuses, getServiceStatus(ctx, k8sClient, nsn))
	}
	reqLogger.Info("New statuses", "statuses", statuses)
	return
}

func getStatefulsetStatus(ctx context.Context, k8sClient client.Client, namespacedName types.NamespacedName) (status operatorv1alpha1.ManagedResourceStatus) {
	reqLogger := logf.FromContext(ctx).WithName("getDeploymentStatus")
	kind := "Statefulset"
	status = operatorv1alpha1.ManagedResourceStatus{
		ObjectName: namespacedName.Name,
		APIVersion: UnknownAPIVersion,
		Namespace:  namespacedName.Namespace,
		Kind:       kind,
		Status:     ResourceNotReadyState,
	}
	statefulset := &appsv1.StatefulSet{}
	err := k8sClient.Get(ctx, namespacedName, statefulset)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Could not find resource for status update", "kind", kind, "name", namespacedName.Name, "namespace", namespacedName.Namespace)
		} else {
			reqLogger.Error(err, "Error reading resource for status update", "kind", kind, "name", namespacedName.Name, "namespace", namespacedName.Namespace)
		}
		return
	}
	status.APIVersion = statefulset.APIVersion
	if statefulset.Status.Replicas == statefulset.Status.ReadyReplicas {
		status.Status = ResourceReadyState
		return
	}
	return
}

func getAllStatefulsetStatus(ctx context.Context, k8sClient client.Client, names []string, namespace string) (statuses []operatorv1alpha1.ManagedResourceStatus) {
	reqLogger := logf.FromContext(ctx).WithName("getAllDeploymentStatus").V(3)
	for _, name := range names {
		nsn := types.NamespacedName{Name: name, Namespace: namespace}
		statuses = append(statuses, getStatefulsetStatus(ctx, k8sClient, nsn))
	}
	reqLogger.Info("New statuses", "statuses", statuses)
	return
}

func getCurrentServiceStatus(ctx context.Context, k8sClient client.Client, mongodb *operatorv1alpha1.MongoDB) (status operatorv1alpha1.ServiceStatus) {
	reqLogger := logf.FromContext(ctx).WithName("getCurrentServiceStatus").V(3)
	type statusRetrieval struct {
		names []string
		f     statusRetrievalFunc
	}

	//
	statusRetrievals := []statusRetrieval{
		{
			names: []string{
				"icp-mongodb",
				"mongodb",
			},
			f: getAllServiceStatus,
		},
		{
			names: []string{
				"icp-mongodb",
			},
			f: getAllStatefulsetStatus,
		},
	}

	kind := "MongoDB"
	status = operatorv1alpha1.ServiceStatus{
		ObjectName:       mongodb.Name,
		Namespace:        mongodb.Namespace,
		APIVersion:       mongodb.APIVersion,
		Kind:             kind,
		ManagedResources: []operatorv1alpha1.ManagedResourceStatus{},
		Status:           ResourceNotReadyState,
	}

	reqLogger.Info("Getting statuses")
	for _, getStatuses := range statusRetrievals {
		status.ManagedResources = append(status.ManagedResources, getStatuses.f(ctx, k8sClient, getStatuses.names, status.Namespace)...)
	}

	for _, managedResourceStatus := range status.ManagedResources {
		if managedResourceStatus.Status == ResourceNotReadyState {
			return
		}
	}
	status.Status = ResourceReadyState
	return
}
