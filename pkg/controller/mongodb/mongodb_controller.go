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

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	operatorv1alpha1 "github.com/IBM/ibm-mongodb-operator/pkg/apis/operator/v1alpha1"
)

var log = logf.Log.WithName("controller_mongodb")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MongoDB Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMongoDB{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mongodb-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MongoDB
	err = c.Watch(&source.Kind{Type: &operatorv1alpha1.MongoDB{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner MongoDB

	return nil
}

// blank assignment to verify that ReconcileMongoDB implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMongoDB{}

// ReconcileMongoDB reconciles a MongoDB object
type ReconcileMongoDB struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a MongoDB object and makes changes based on the state read
// and what is in the MongoDB.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMongoDB) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MongoDB")

	// Fetch the MongoDB instance
	instance := &operatorv1alpha1.MongoDB{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	log.Info("creating mongodb service")
	if err := r.createFromYaml(instance, []byte(service)); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("creating mongodb icp service")
	if err := r.createFromYaml(instance, []byte(icpService)); err != nil {
		return reconcile.Result{}, err
	}

	metadatalabel := map[string]string{"app.kubernetes.io/name": "icp-mongodb", "app.kubernetes.io/component": "database",
		"app.kubernetes.io/managed-by": "operator", "app.kubernetes.io/instance": "icp-mongodb", "release": "mongodb"}

	log.Info("creating icp mongodb config map")

	if err := r.createFromYaml(instance, []byte(mongodbConfigMap)); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("creating icp mongodb init config map")

	if err := r.createFromYaml(instance, []byte(initConfigMap)); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("creating icp mongodb install config map")

	if err := r.createFromYaml(instance, []byte(installConfigMap)); err != nil {
		return reconcile.Result{}, err
	}

	var pass, user string
	if instance.Spec.MongoDBPass == "" {
		pass = "admin"
	} else {
		pass = instance.Spec.MongoDBPass
	}

	if instance.Spec.MongoDBUser == "" {
		user = "admin"
	} else {
		user = instance.Spec.MongoDBUser
	}

	mongodbAdmin := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "icp-mongodb",
			},
			Name:      "icp-mongodb-admin",
			Namespace: instance.GetNamespace(),
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"user":     user,
			"password": pass,
		},
	}

	// Set CommonServiceConfig instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, mongodbAdmin, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("creating icp mongodb admin secret")
	if err = r.client.Create(context.TODO(), mongodbAdmin); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	}

	mongodbMetric := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    metadatalabel,
			Name:      "icp-mongodb-metrics",
			Namespace: instance.GetNamespace(),
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"user":     "metrics",
			"password": "icpmetrics",
		},
	}

	// Set CommonServiceConfig instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, mongodbMetric, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("creating icp mongodb metric secret")
	if err = r.client.Create(context.TODO(), mongodbMetric); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	}

	keyfileSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    metadatalabel,
			Name:      "icp-mongodb-keyfile",
			Namespace: instance.GetNamespace(),
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"key.txt": "icptest",
		},
	}

	// Set CommonServiceConfig instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, keyfileSecret, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("creating icp mongodb keyfile secret")
	if err = r.client.Create(context.TODO(), keyfileSecret); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	}

	var storageclass string

	if instance.Status.StorageClass == "" {
		if instance.Spec.StorageClass == "" {
			// TODO: weird because the storage class on OCP is opened for all
			// Need to deploy an OCP cluster on AWS to verify
			storageclass, err = r.getstorageclass()
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			storageclass = instance.Spec.StorageClass
		}
	} else {
		if instance.Spec.StorageClass != "" && instance.Spec.StorageClass != instance.Status.StorageClass {
			log.Info("You need to delete the monogodb cr before switch the storage class. Please note that this will lose all your datamake")
		}
		storageclass = instance.Status.StorageClass
	}

	fmt.Println(storageclass)

	stsData := struct {
		Replicas     int
		ImageRepo    string
		StorageClass string
	}{
		Replicas:     instance.Spec.Replicas,
		ImageRepo:    instance.Spec.ImageRegistry,
		StorageClass: storageclass,
	}

	var stsYaml bytes.Buffer
	t := template.Must(template.New("statefulset").Parse(statefulset))
	if err := t.Execute(&stsYaml, stsData); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("creating mongodb statefulset")
	if err := r.createFromYaml(instance, stsYaml.Bytes()); err != nil {
		return reconcile.Result{}, err
	}

	instance.Status.StorageClass = storageclass
	if err := r.client.Status().Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, err
	}

	// hardcode cluster-ca-cert here
	// TODO: later will copy it from cert-manager namespace
	log.Info("creating cluster-ca-cert")
	if err := r.createFromYaml(instance, []byte(clusterCertYaml)); err != nil {
		log.Error(err, "create cluster-ca-cert fail")
	}
	log.Info("creating icp-mongodb-client-cert")
	if err := r.createFromYaml(instance, []byte(mongoCertYaml)); err != nil {
		log.Error(err, "create icp-mongodb-client-cert fail")
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileMongoDB) createFromYaml(instance *operatorv1alpha1.MongoDB, yamlContent []byte) error {
	obj := &unstructured.Unstructured{}
	jsonSpec, err := yaml.YAMLToJSON(yamlContent)
	if err != nil {
		return fmt.Errorf("could not convert yaml to json: %v", err)
	}

	if err := obj.UnmarshalJSON(jsonSpec); err != nil {
		return fmt.Errorf("could not unmarshal resource: %v", err)
	}

	// Set CommonServiceConfig instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, obj, r.scheme); err != nil {
		return err
	}

	err = r.client.Create(context.TODO(), obj)
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("could not Create resource: %v", err)
	}

	return nil
}

func (r *ReconcileMongoDB) getstorageclass() (string, error) {
	scList := &storagev1.StorageClassList{}
	err := r.client.List(context.TODO(), scList)
	if err != nil {
		return "", err
	}
	if len(scList.Items) == 0 {
		return "", fmt.Errorf("could not find storage class in the cluster")
	}

	var defaultSC []string
	var nonDefaultSC []string

	for _, sc := range scList.Items {
		if sc.Provisioner == "kubernetes.io/no-provisioner" {
			continue
		}
		if sc.ObjectMeta.GetAnnotations()["storageclass.kubernetes.io/is-default-class"] == "true" {
			defaultSC = append(defaultSC, sc.GetName())
			continue
		}
		nonDefaultSC = append(nonDefaultSC, sc.GetName())
	}

	if len(defaultSC) != 0 {
		return defaultSC[0], nil
	}

	if len(nonDefaultSC) != 0 {
		return nonDefaultSC[0], nil
	}

	return "", fmt.Errorf("could not find dynamic provisioner storage class in the cluster")
}
