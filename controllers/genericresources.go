//
// Copyright 2023 IBM Corporation
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

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func DeleteGenericResource(ctx context.Context, name string, namespace string, group string, version string, resource string) error {
	reqLogger := logf.FromContext(ctx).WithName("DeleteGenericResource")

	config := ctrl.GetConfigOrDie()
	dynamic := dynamic.NewForConfigOrDie(config)

	resourceID := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	_, err := dynamic.Resource(resourceID).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Unstructured resource not found")
			return nil
		}
		reqLogger.Error(err, "Failed reading unstructured resource")
		return err
	}

	// Delete ibm-mongodb-request operandrequest if found
	err = dynamic.Resource(resourceID).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		reqLogger.Error(err, "Failed to delete unstructured resource")
		return err
	}

	reqLogger.Info("Deleted unstructured resource ibm-mongodb-request operandrequest")
	return nil
}
