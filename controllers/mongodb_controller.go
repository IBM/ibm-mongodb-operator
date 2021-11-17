/*
# Copyright 2021 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	mongodbv1alpha1 "github.com/IBM/ibm-mongodb-operator/api/v1alpha1"
)

// MongoDBReconciler reconciles a MongoDB object
type MongoDBReconciler struct {
	Client client.Client
	Reader client.Reader
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//
const mongodbOperatorURI = `mongodbs.operator.ibm.com`
const defaultPVCSize = `20Gi`

// MongoDB StatefulSet Data
type mongoDBStatefulSetData struct {
	Replicas       int
	ImageRepo      string
	StorageClass   string
	InitImage      string
	BootstrapImage string
	MetricsImage   string
	CPULimit       string
	CPURequest     string
	MemoryLimit    string
	MemoryRequest  string
	NamespaceName  string
	StsLabels      map[string]string
	PodLabels      map[string]string
	PVCSize        string
	UserId         int
}

// +kubebuilder:rbac:groups=mongodb.operator.ibm.com,namespace=ibm-common-services,resources=mongodbs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mongodb.operator.ibm.com,namespace=ibm-common-services,resources=mongodbs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,namespace=ibm-common-services,resources=services;services/finalizers;serviceaccounts;endpoints;persistentvolumeclaims;events;configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,namespace=ibm-common-services,resources=deployments;daemonsets;replicasets;statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,namespace=ibm-common-services,resources=servicemonitors,verbs=get;create
// +kubebuilder:rbac:groups=apps,namespace=ibm-common-services,resourceNames=ibm-mongodb-operator,resources=deployments/finalizers,verbs=update
// +kubebuilder:rbac:groups=operator.ibm.com,namespace=ibm-common-services,resources=mongodbs;mongodbs/finalizers;mongodbs/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=certmanager.k8s.io,namespace=ibm-common-services,resources=certificates;certificaterequests;orders;challenges;issuers,verbs=get;list;watch;create;update;patch;delete

func (r *MongoDBReconciler) Reconcile(request ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("mongodb", request.NamespacedName)

	// Fetch the MongoDB instance
	instance := &mongodbv1alpha1.MongoDB{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, instance)
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

	r.Log.Info("creating mongodb service account")
	if err := r.createFromYaml(instance, []byte(mongoSA)); err != nil {
		return reconcile.Result{}, err
	}

	r.Log.Info("creating mongodb service")
	if err := r.createFromYaml(instance, []byte(service)); err != nil {
		return reconcile.Result{}, err
	}

	r.Log.Info("creating mongodb icp service")
	if err := r.createFromYaml(instance, []byte(icpService)); err != nil {
		return reconcile.Result{}, err
	}

	metadatalabel := map[string]string{"app.kubernetes.io/name": "icp-mongodb", "app.kubernetes.io/component": "database",
		"app.kubernetes.io/managed-by": "operator", "app.kubernetes.io/instance": "icp-mongodb", "release": "mongodb"}

	r.Log.Info("creating icp mongodb config map")
	//Calculate MongoDB cache Size
	var cacheSize float64
	var cacheSizeGB float64
	if instance.Spec.Resources.Limits.Memory().String() != "0" {
		ramMB := instance.Spec.Resources.Limits.Memory().ScaledValue(resource.Mega)
		// Cache Size is 40 percent of RAM
		cacheSize = float64(ramMB) * 0.4
		// Convert to gig
		cacheSizeGB = cacheSize / 1000.0
		// Round to fit config
		cacheSizeGB = math.Floor(cacheSizeGB*100) / 100
	} else {
		//default value is 5Gi
		cacheSizeGB = 2.0
	}

	monogdbConfigmapData := struct {
		CacheSize float64
	}{
		CacheSize: cacheSizeGB,
	}
	// TO DO -- convert configmap to take option.
	var mongodbConfigYaml bytes.Buffer
	tc := template.Must(template.New("mongodbconfigmap").Parse(mongodbConfigMap))
	if err := tc.Execute(&mongodbConfigYaml, monogdbConfigmapData); err != nil {
		return reconcile.Result{}, err
	}

	r.Log.Info("creating or updating mongodb configmap")
	if err := r.createUpdateFromYaml(instance, mongodbConfigYaml.Bytes()); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.createFromYaml(instance, []byte(mongodbConfigMap)); err != nil {
		return reconcile.Result{}, err
	}

	r.Log.Info("creating or updating icp mongodb init config map")

	if err := r.createUpdateFromYaml(instance, []byte(initConfigMap)); err != nil {
		return reconcile.Result{}, err
	}

	r.Log.Info("creating icp mongodb install config map")

	if err := r.createFromYaml(instance, []byte(installConfigMap)); err != nil {
		return reconcile.Result{}, err
	}

	// Create admin user and password as random string
	// TODO: allow user to give a Secret
	var pass, user string
	user = createRandomAlphaNumeric(8)
	pass = createRandomAlphaNumeric(13)
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
	// if err := controllerutil.SetControllerReference(instance, mongodbAdmin, r.scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }

	r.Log.Info("creating icp mongodb admin secret")
	if err = r.Client.Create(context.TODO(), mongodbAdmin); err != nil && !errors.IsAlreadyExists(err) {
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
	if err := controllerutil.SetControllerReference(instance, mongodbMetric, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	r.Log.Info("creating icp mongodb metric secret")
	if err = r.Client.Create(context.TODO(), mongodbMetric); err != nil && !errors.IsAlreadyExists(err) {
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
	if err := controllerutil.SetControllerReference(instance, keyfileSecret, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	r.Log.Info("creating icp mongodb keyfile secret")
	if err = r.Client.Create(context.TODO(), keyfileSecret); err != nil && !errors.IsAlreadyExists(err) {
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
			r.Log.Info("You need to delete the monogodb cr before switch the storage class. Please note that this will lose all your datamake")
		}
		storageclass = instance.Status.StorageClass
	}

	// Default values
	cpuRequest := "2000m"
	memoryRequest := "5Gi"
	cpuLimit := "2000m"
	memoryLimit := "5Gi"

	// Check cpu request values and default if not there
	if instance.Spec.Resources.Requests.Cpu().String() != "0" {
		cpuRequest = instance.Spec.Resources.Requests.Cpu().String()
	}

	// Check memory request values and default if not there
	if instance.Spec.Resources.Requests.Memory().String() != "0" {
		memoryRequest = instance.Spec.Resources.Requests.Memory().String()
	}

	// Check cpu limit values and default if not there
	if instance.Spec.Resources.Limits.Cpu().String() != "0" {
		cpuLimit = instance.Spec.Resources.Limits.Cpu().String()
	}

	// Check memory limit values and default if not there
	if instance.Spec.Resources.Limits.Memory().String() != "0" {
		memoryLimit = instance.Spec.Resources.Limits.Memory().String()
	}

	// Default values
	PVCSizeRequest := defaultPVCSize

	// If PVC already exist and the value does not match the PVCSizeRequest then log information that it cannot be changed.
	pvc := &corev1.PersistentVolumeClaim{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: "mongodbdir-icp-mongodb-0", Namespace: instance.Namespace}, pvc)
	if err == nil {
		PVCSizeRequest = pvc.Spec.Resources.Requests.Storage().String()
		if instance.Spec.PVC.Resources.Requests.Storage().String() != "0" {
			if (PVCSizeRequest != instance.Spec.PVC.Resources.Requests.Storage().String()) && (instance.Spec.PVC.Resources.Requests.Storage().String() != defaultPVCSize) {
				r.Log.Info("mongoDB Persistent Volume Claim already exists, it's size is immutable, ignoring requested storage size for the PVC")
			}
		} else {
			if PVCSizeRequest != defaultPVCSize {
				r.Log.Info("mongoDB Persistent Volume Claim already exists, it's size is immutable.")
				r.Log.Info("the PVC storage request is not set to the current default nor is it specified in the Custom Resource")
			}
		}
	} else if errors.IsNotFound(err) {
		// Check PVC size request values and default if not there
		if instance.Spec.PVC.Resources.Requests.Storage().String() != "0" {
			PVCSizeRequest = instance.Spec.PVC.Resources.Requests.Storage().String()
		}
	}

	// Select User to use
	cppConfig := &corev1.ConfigMap{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: "ibm-cpp-config", Namespace: instance.Namespace}, cppConfig)
	if err != nil {
		return reconcile.Result{}, err
	}

	uid := 0
	if cppConfig.Data["kubernetes_cluster_type"] == "cncf" {
		uid = 1000
	}

	// Check if statefulset already exists
	sts := &appsv1.StatefulSet{}
	var stsLabels map[string]string
	var podLabels map[string]string

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: "icp-mongodb", Namespace: instance.Namespace}, sts)
	if err == nil {
		r.Log.Info("succeeded to get statefulset check")
		stsLabels = sts.ObjectMeta.Labels
		podLabels = sts.Spec.Template.ObjectMeta.Labels
	} else if errors.IsNotFound(err) {
		r.Log.Info("statefulset not found for labels")
		constStsLabels := make(map[string]string)
		constStsLabels["app"] = "icp-mongodb"
		constStsLabels["release"] = "mongodb"
		constStsLabels["app.kubernetes.io/instance"] = mongodbOperatorURI
		constStsLabels["app.kubernetes.io/managed-by"] = mongodbOperatorURI
		constStsLabels["app.kubernetes.io/name"] = mongodbOperatorURI
		stsLabels = constStsLabels
		constPodLabels := make(map[string]string)
		constPodLabels["app.kubernetes.io/instance"] = "common-mongodb"
		constPodLabels["app"] = "icp-mongodb"
		constPodLabels["release"] = "mongodb"
		podLabels = constPodLabels
	} else {
		return reconcile.Result{}, err
	}

	//Set Replicas
	//Get current number of replicas in cluster based on number of PVCs
	pvcs := &corev1.PersistentVolumeClaimList{}
	err = r.Client.List(context.TODO(), pvcs, &client.ListOptions{
		Namespace: instance.Namespace,
	})
	pvcCount := 0
	if err == nil {
		// loop items in pvcs and count mongodbdir
		for _, pvc := range pvcs.Items {
			if strings.Contains(pvc.ObjectMeta.Name, "mongodbdir-icp-mongodb") {
				pvcCount++
				r.Log.Info("Found PVC" + pvc.ObjectMeta.Name)
			}
		}
	} else {
		return reconcile.Result{}, err
	}

	//check pvc count with replicas
	//if pvcCount is greater than the replicas, then at one time there must have been more replicas
	replicas := instance.Spec.Replicas
	if pvcCount > replicas {
		replicas = pvcCount
		r.Log.Info("Ignoring Replica spec, there are more mongodbdir-icp-mongodb PVCs than the current relica request.")
		r.Log.Info("PVC count: " + strconv.Itoa(pvcCount))
	}
	stsData := mongoDBStatefulSetData{
		Replicas:       replicas,
		ImageRepo:      instance.Spec.ImageRegistry,
		StorageClass:   storageclass,
		InitImage:      os.Getenv("IBM_MONGODB_INSTALL_IMAGE"),
		BootstrapImage: os.Getenv("IBM_MONGODB_IMAGE"),
		MetricsImage:   os.Getenv("IBM_MONGODB_EXPORTER_IMAGE"),
		CPULimit:       cpuLimit,
		CPURequest:     cpuRequest,
		MemoryLimit:    memoryLimit,
		MemoryRequest:  memoryRequest,
		NamespaceName:  instance.Namespace,
		StsLabels:      stsLabels,
		PodLabels:      podLabels,
		PVCSize:        PVCSizeRequest,
		UserId:         uid,
	}

	var stsYaml bytes.Buffer
	t := template.Must(template.New("statefulset").Parse(statefulset))
	if err := t.Execute(&stsYaml, stsData); err != nil {
		return reconcile.Result{}, err
	}

	r.Log.Info("creating mongodb statefulset")
	if err := r.createUpdateFromYaml(instance, stsYaml.Bytes()); err != nil {
		return reconcile.Result{}, err
	}

	instance.Status.StorageClass = storageclass
	if err := r.Client.Status().Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, err
	}

	// sign certificate
	r.Log.Info("creating root-ca-cert")
	if err := r.createFromYaml(instance, []byte(godIssuerYaml)); err != nil {
		r.Log.Error(err, "create god-issuer fail")
		return reconcile.Result{}, err
	}
	r.Log.Info("creating root-ca-cert")
	if err := r.createFromYaml(instance, []byte(rootCertYaml)); err != nil {
		r.Log.Error(err, "create root-ca-cert fail")
		return reconcile.Result{}, err
	}
	r.Log.Info("creating root-issuer")
	if err := r.createFromYaml(instance, []byte(rootIssuerYaml)); err != nil {
		r.Log.Error(err, "create root-issuer fail")
		return reconcile.Result{}, err
	}
	r.Log.Info("creating icp-mongodb-client-cert")
	if err := r.createFromYaml(instance, []byte(clientCertYaml)); err != nil {
		r.Log.Error(err, "create icp-mongodb-client-cert fail")
		return reconcile.Result{}, err
	}

	// Get the StatefulSet
	sts = &appsv1.StatefulSet{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{Name: "icp-mongodb", Namespace: instance.Namespace}, sts); err != nil {
		return reconcile.Result{}, err
	}

	// Add controller on PVC
	if err = r.addControlleronPVC(instance, sts); err != nil {
		return reconcile.Result{}, err
	}

	if sts.Status.UpdatedReplicas != sts.Status.Replicas || sts.Status.UpdatedReplicas != sts.Status.ReadyReplicas {
		r.Log.Info("Waiting Mongodb to be ready ...")
		return reconcile.Result{Requeue: true, RequeueAfter: time.Minute}, nil
	}
	r.Log.Info("Mongodb is ready")

	return ctrl.Result{}, nil
}

// Move to separate file begin

func (r *MongoDBReconciler) createFromYaml(instance *mongodbv1alpha1.MongoDB, yamlContent []byte) error {
	obj := &unstructured.Unstructured{}
	jsonSpec, err := yaml.YAMLToJSON(yamlContent)
	if err != nil {
		return fmt.Errorf("could not convert yaml to json: %v", err)
	}

	if err := obj.UnmarshalJSON(jsonSpec); err != nil {
		return fmt.Errorf("could not unmarshal resource: %v", err)
	}

	obj.SetNamespace(instance.Namespace)

	// Set CommonServiceConfig instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, obj, r.Scheme); err != nil {
		return err
	}

	err = r.Client.Create(context.TODO(), obj)
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("could not Create resource: %v", err)
	}

	return nil
}

func (r *MongoDBReconciler) createUpdateFromYaml(instance *mongodbv1alpha1.MongoDB, yamlContent []byte) error {
	obj := &unstructured.Unstructured{}
	jsonSpec, err := yaml.YAMLToJSON(yamlContent)
	if err != nil {
		return fmt.Errorf("could not convert yaml to json: %v", err)
	}

	if err := obj.UnmarshalJSON(jsonSpec); err != nil {
		return fmt.Errorf("could not unmarshal resource: %v", err)
	}

	obj.SetNamespace(instance.Namespace)

	// Set CommonServiceConfig instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, obj, r.Scheme); err != nil {
		return err
	}

	err = r.Client.Create(context.TODO(), obj)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			if err := r.Client.Update(context.TODO(), obj); err != nil {
				return fmt.Errorf("could not Update resource: %v", err)
			}
			return nil
		}
		return fmt.Errorf("could not Create resource: %v", err)
	}

	return nil
}

func (r *MongoDBReconciler) getstorageclass() (string, error) {
	scList := &storagev1.StorageClassList{}
	err := r.Reader.List(context.TODO(), scList)
	if err != nil {
		return "", err
	}
	if len(scList.Items) == 0 {
		return "", fmt.Errorf("could not find storage class in the cluster")
	}

	var defaultSC []string
	var nonDefaultSC []string

	for _, sc := range scList.Items {
		if sc.ObjectMeta.GetAnnotations()["storageclass.kubernetes.io/is-default-class"] == "true" {
			defaultSC = append(defaultSC, sc.GetName())
			continue
		}
		if sc.Provisioner == "kubernetes.io/no-provisioner" {
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

	return "", fmt.Errorf("could not find dynamic provisioner storage class in the cluster nor is there a default storage class")
}

func (r *MongoDBReconciler) addControlleronPVC(instance *mongodbv1alpha1.MongoDB, sts *appsv1.StatefulSet) error {
	// Fetch the list of the PersistentVolumeClaim generated by the StatefulSet
	pvcList := &corev1.PersistentVolumeClaimList{}
	err := r.Client.List(context.TODO(), pvcList, &client.ListOptions{
		Namespace:     instance.Namespace,
		LabelSelector: labels.SelectorFromSet(sts.ObjectMeta.Labels),
	})

	if err != nil {
		return err
	}

	for _, pvc := range pvcList.Items {
		if pvc.ObjectMeta.OwnerReferences == nil {
			if err := controllerutil.SetControllerReference(instance, &pvc, r.Scheme); err != nil {
				return err
			}
			if err = r.Client.Update(context.TODO(), &pvc); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create Random String
func createRandomAlphaNumeric(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	byteString := make([]byte, length)
	for i := range byteString {
		byteString[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(byteString)
}

// Move to separate file?

func (r *MongoDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mongodbv1alpha1.MongoDB{}).
		Owns(&appsv1.StatefulSet{}).Owns(&corev1.ConfigMap{}).Owns(&corev1.ServiceAccount{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
