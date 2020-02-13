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
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorv1alpha1.MongoDB{},
	})
	if err != nil {
		return err
	}

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

	stsData := struct {
		Replicas     int
		ImageRepo    string
		StorageClass string
	}{
		Replicas:     instance.Spec.Replicas,
		ImageRepo:    instance.Spec.ImageRegistry,
		StorageClass: instance.Spec.StorageClass,
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

	metadatalabel := map[string]string{"app.kubernetes.io/name": "icp-mongodb", "app.kubernetes.io/component": "database",
		"app.kubernetes.io/managed-by": "operator", "app.kubernetes.io/instance": "icp-mongodb", "release": "mongodb"}
	icpMongodbConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    metadatalabel,
			Name:      "icp-mongodb",
			Namespace: instance.GetNamespace(),
		},
		Data: map[string]string{
			"mongod.conf": `|
    storage:
      dbPath: /data/db
    net:
      bindIpAll: true
      port: 27017
      ssl:
        mode: requireSSL
        CAFile: /data/configdb/tls.crt
        PEMKeyFile: /work-dir/mongo.pem
    replication:
      replSetName: rs0
    # Uncomment for TLS support or keyfile access control without TLS
    security:
      authorization: enabled
      keyFile: /data/configdb/key.txt
			`,
		},
	}

	// Set CommonServiceConfig instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, icpMongodbConfigMap, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("creating icp mongodb config map")

	if err = r.client.Create(context.TODO(), icpMongodbConfigMap); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	}

	icpMongodbinitConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    metadatalabel,
			Name:      "icp-mongodb-init",
			Namespace: instance.GetNamespace(),
		},
		Data: map[string]string{
			"on-start.sh": `|
			#!/bin/bash
    
			port=27017
			replica_set=$REPLICA_SET
			script_name=${0##*/}
			credentials_file=/work-dir/credentials.txt
			config_dir=/data/configdb
			
			function log() {
				local msg="$1"
				local timestamp=$(date --iso-8601=ns)
				echo "[$timestamp] [$script_name] $msg"
				echo "[$timestamp] [$script_name] $msg" >> /work-dir/log.txt
			}
			
			if [[ "$AUTH" == "true" ]]; then
			
				if [ !  -f "$credentials_file" ]; then
					log "Creds File Not found!"
					echo $ADMIN_USER > $credentials_file
					echo $ADMIN_PASSWORD >> $credentials_file
				fi
				admin_user=$(head -n 1 $credentials_file)
				admin_password=$(tail -n 1 $credentials_file)
				admin_auth=(-u "$admin_user" -p "$admin_password")
				if [[ "$METRICS" == "true" ]]; then
					metrics_user="$METRICS_USER"
					metrics_password="$METRICS_PASSWORD"
				fi
			fi
			
			function shutdown_mongo() {
				if [[ $# -eq 1 ]]; then
					args="timeoutSecs: $1"
				else
					args='force: true'
				fi
				log "Shutting down MongoDB ($args)..."
				mongo admin "${admin_auth[@]}" "${ssl_args[@]}" --eval "db.shutdownServer({$args})"
			}
			
			#Check if Password has change and updated in mongo , if so update Creds
			function update_creds_if_changed() {
			  if [ "$admin_password" != "$ADMIN_PASSWORD" ]; then
				  passwd_changed=true
				  log "password has changed = $passwd_changed"
				  log "checking if passwd  updated in mongo"
				  mongo admin  "${ssl_args[@]}" --eval "db.auth({user: '$admin_user', pwd: '$ADMIN_PASSWORD'})" | grep "Authentication failed"
				  if [[ $? -eq 1 ]]; then
					log "New Password worked, update creds"
					echo $ADMIN_USER > $credentials_file
					echo $ADMIN_PASSWORD >> $credentials_file
					admin_password=$ADMIN_PASSWORD
					admin_auth=(-u "$admin_user" -p "$admin_password")
					passwd_updated=true
				  fi
			  fi
			}
			
			function update_mongo_password_if_changed() {
			  log "checking if mongo passwd needs to be  updated"
			  if [[ "$passwd_changed" == "true" ]] && [[ "$passwd_updated" != "true" ]]; then
				log "Updating to new password "
				if [[ $# -eq 1 ]]; then
					mhost="--host $1"
				else
					mhost=""
				fi
			
				log "host for password upd ($mhost)"
				mongo admin $mhost "${admin_auth[@]}" "${ssl_args[@]}" --eval "db.changeUserPassword('$admin_user', '$ADMIN_PASSWORD')" >> /work-dir/log.txt 2>&1
				sleep 10
				log "mongo passwd change attempted; check and update creds file if successfully"
				update_creds_if_changed
			  fi
			}
			
			
			
			my_hostname=$(hostname)
			log "Bootstrapping MongoDB replica set member: $my_hostname"
			
			log "Reading standard input..."
			while read -ra line; do
				log "line is  ${line}"
				if [[ "${line}" == *"${my_hostname}"* ]]; then
					service_name="$line"
				fi
				peers=("${peers[@]}" "$line")
			done
			
			# Move into /work-dir
			pushd /work-dir
			pwd >> /work-dir/log.txt
			ls -l  >> /work-dir/log.txt
			
			# Generate the ca cert
			ca_crt=$config_dir/tls.crt
			if [ -f $ca_crt  ]; then
				log "Generating certificate"
				ca_key=$config_dir/tls.key
				pem=/work-dir/mongo.pem
				ssl_args=(--ssl --sslCAFile $ca_crt --sslPEMKeyFile $pem)
			
				echo "ca stuff" >> /work-dir/log.txt
				cat $ca_crt >> /work-dir/log.txt
				cat $ca_key >> /work-dir/log.txt
			
			cat >openssl.cnf <<EOL
			[req]
			req_extensions = v3_req
			distinguished_name = req_distinguished_name
			[req_distinguished_name]
			[ v3_req ]
			basicConstraints = CA:FALSE
			keyUsage = nonRepudiation, digitalSignature, keyEncipherment
			subjectAltName = @alt_names
			[alt_names]
			DNS.1 = $(echo -n "$my_hostname" | sed s/-[0-9]*$//)
			DNS.2 = $my_hostname
			DNS.3 = $service_name
			DNS.4 = localhost
			DNS.5 = 127.0.0.1
			DNS.6 = mongodb
			EOL
			
				# Generate the certs
				echo "cnf stuff" >> /work-dir/log.txt
				cat openssl.cnf >> /work-dir/log.txt
				echo "genrsa " >> /work-dir/log.txt
				openssl genrsa -out mongo.key 2048 >> /work-dir/log.txt 2>&1
			
				echo "req " >> /work-dir/log.txt
				openssl req -new -key mongo.key -out mongo.csr -subj "/CN=$my_hostname" -config openssl.cnf >> /work-dir/log.txt 2>&1
			
				echo "x509 " >> /work-dir/log.txt
				openssl x509 -req -in mongo.csr \
					-CA $ca_crt -CAkey $ca_key -CAcreateserial \
					-out mongo.crt -days 3650 -extensions v3_req -extfile openssl.cnf >> /work-dir/log.txt 2>&1
			
				echo "mongo stuff" >> /work-dir/log.txt
				cat mongo.csr >> /work-dir/log.txt
			
				rm mongo.csr
			
				echo "mongo key" >> /work-dir/log.txt
				cat mongo.key >> /work-dir/log.txt
				echo "mongo crt" >> /work-dir/log.txt
				cat mongo.crt >> /work-dir/log.txt
				cat mongo.crt mongo.key > $pem
				rm mongo.key mongo.crt
			fi
			
			
			log "Peers: ${peers[@]}"
			
			log "Starting a MongoDB instance..."
			mongod --config $config_dir/mongod.conf >> /work-dir/log.txt 2>&1 &
			pid=$!
			trap shutdown_mongo EXIT
			
			
			log "Waiting for MongoDB to be ready..."
			until [[ $(mongo "${ssl_args[@]}" --quiet --eval "db.adminCommand('ping').ok") == "1" ]]; do
				log "Retrying..."
				sleep 2
			done
			
			log "Initialized."
			
			if [[ "$AUTH" == "true" ]]; then
				update_creds_if_changed
			fi
			
			iter_counter=0
			while [  $iter_counter -lt 5 ]; do
			  log "primary check, iter_counter is $iter_counter"
			  # try to find a master and add yourself to its replica set.
			  for peer in "${peers[@]}"; do
				  log "Checking if ${peer} is primary"
				  mongo admin --host "${peer}" --ipv6 "${admin_auth[@]}" "${ssl_args[@]}" --quiet --eval "rs.status()"  >> log.txt
			
				  # Check rs.status() first since it could be in primary catch up mode which db.isMaster() doesn't show
				  if [[ $(mongo admin --host "${peer}" --ipv6 "${admin_auth[@]}" "${ssl_args[@]}" --quiet --eval "rs.status().myState") == "1" ]]; then
					  log "Found master ${peer}, wait while its in primary catch up mode "
					  until [[ $(mongo admin --host "${peer}" --ipv6 "${admin_auth[@]}" "${ssl_args[@]}" --quiet --eval "db.isMaster().ismaster") == "true" ]]; do
						  sleep 1
					  done
					  primary="${peer}"
					  log "Found primary: ${primary}"
					  break
				  fi
			  done
			
			  if [[ -z "${primary}" ]]  && [[ ${#peers[@]} -gt 1 ]] && (mongo "${ssl_args[@]}" --eval "rs.status()" | grep "no replset config has been received"); then
				log "waiting before creating a new replicaset, to avoid conflicts with other replicas"
				sleep 30
			  else
				break
			  fi
			
			  let iter_counter=iter_counter+1
			done
			
			
			if [[ "${primary}" = "${service_name}" ]]; then
				log "This replica is already PRIMARY"
			
			elif [[ -n "${primary}" ]]; then
			
				if [[ $(mongo admin --host "${primary}" --ipv6 "${admin_auth[@]}" "${ssl_args[@]}" --quiet --eval "rs.conf().members.findIndex(m => m.host == '${service_name}:${port}')") == "-1" ]]; then
				  log "Adding myself (${service_name}) to replica set..."
				  if (mongo admin --host "${primary}" --ipv6 "${admin_auth[@]}" "${ssl_args[@]}" --eval "rs.add('${service_name}')" | grep 'Quorum check failed'); then
					  log 'Quorum check failed, unable to join replicaset. Exiting.'
					  exit 1
				  fi
				fi
				log "Done,  Added myself to replica set."
			
				sleep 3
				log 'Waiting for replica to reach SECONDARY state...'
				until printf '.'  && [[ $(mongo admin "${admin_auth[@]}" "${ssl_args[@]}" --quiet --eval "rs.status().myState") == '2' ]]; do
					sleep 1
				done
				log '✓ Replica reached SECONDARY state.'
			
			elif (mongo "${ssl_args[@]}" --eval "rs.status()" | grep "no replset config has been received"); then
			
				log "Initiating a new replica set with myself ($service_name)..."
			
				mongo "${ssl_args[@]}" --eval "rs.initiate({'_id': '$replica_set', 'members': [{'_id': 0, 'host': '$service_name'}]})"
				mongo "${ssl_args[@]}" --eval "rs.status()"
			
				sleep 3
			
				log 'Waiting for replica to reach PRIMARY state...'
			
				log ' Waiting for rs.status state to become 1'
				until printf '.'  && [[ $(mongo "${ssl_args[@]}" --quiet --eval "rs.status().myState") == '1' ]]; do
					sleep 1
				done
			
				log ' Waiting for master to complete primary catchup mode'
				until [[ $(mongo  "${ssl_args[@]}" --quiet --eval "db.isMaster().ismaster") == "true" ]]; do
					sleep 1
				done
			
				primary="${service_name}"
				log '✓ Replica reached PRIMARY state.'
			
			
				if [[ "$AUTH" == "true" ]]; then
					# sleep a little while just to be sure the initiation of the replica set has fully
					# finished and we can create the user
					sleep 3
			
					log "Creating admin user..."
					mongo admin "${ssl_args[@]}" --eval "db.createUser({user: '$admin_user', pwd: '$admin_password', roles: [{role: 'root', db: 'admin'}]})"
				fi
			
				log "Done initiating replicaset."
			
			fi
			
			log "Primary: ${primary}"
			
			if [[  -n "${primary}"   && "$AUTH" == "true" ]]; then
				# you r master and passwd has changed.. then update passwd
				update_mongo_password_if_changed $primary
			
				if [[ "$METRICS" == "true" ]]; then
					log "Checking if metrics user is already created ..."
					metric_user_count=$(mongo admin --host "${primary}" "${admin_auth[@]}" "${ssl_args[@]}" --eval "db.system.users.find({user: '${metrics_user}'}).count()" --quiet)
					log "User count is ${metric_user_count} "
					if [[ "${metric_user_count}" == "0" ]]; then
						log "Creating clusterMonitor user... user - ${metrics_user}  "
						mongo admin --host "${primary}" "${admin_auth[@]}" "${ssl_args[@]}" \
						--eval "db.createUser({user: '${metrics_user}', pwd: '${metrics_password}', roles: [{role: 'clusterMonitor', db: 'admin'}, {role: 'read', db: 'local'}]})"
						log "User creation return code is $? "
						metric_user_count=$(mongo admin --host "${primary}" "${admin_auth[@]}" "${ssl_args[@]}" --eval "db.system.users.find({user: '${metrics_user}'}).count()" --quiet)
						log "User count now is ${metric_user_count} "
					fi
				fi
			fi
			
			log "MongoDB bootstrap complete"
			exit 0
		
			`,
		},
	}

	// Set CommonServiceConfig instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, icpMongodbinitConfigMap, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("creating icp mongodb init config map")

	if err = r.client.Create(context.TODO(), icpMongodbinitConfigMap); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	}

	icpMongodbInstallConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    metadatalabel,
			Name:      "icp-mongodb-install",
			Namespace: instance.GetNamespace(),
		},
		Data: map[string]string{
			"install.sh": `|
			#!/bin/bash
    
			# Copyright 2016 The Kubernetes Authors. All rights reserved.
			#
			# Licensed under the Apache License, Version 2.0 (the "License");
			# you may not use this file except in compliance with the License.
			# You may obtain a copy of the License at
			#
			#     http://www.apache.org/licenses/LICENSE-2.0
			#
			# Unless required by applicable law or agreed to in writing, software
			# distributed under the License is distributed on an "AS IS" BASIS,
			# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
			# See the License for the specific language governing permissions and
			# limitations under the License.
			
			# This volume is assumed to exist and is shared with the peer-finder
			# init container. It contains on-start/change configuration scripts.
			WORKDIR_VOLUME="/work-dir"
			CONFIGDIR_VOLUME="/data/configdb"
			
			for i in "$@"
			do
			case $i in
				-c=*|--config-dir=*)
				CONFIGDIR_VOLUME="${i#*=}"
				shift
				;;
				-w=*|--work-dir=*)
				WORKDIR_VOLUME="${i#*=}"
				shift
				;;
				*)
				# unknown option
				;;
			esac
			done
			
			echo installing config scripts into "${WORKDIR_VOLUME}"
			mkdir -p "${WORKDIR_VOLUME}"
			cp /peer-finder "${WORKDIR_VOLUME}"/
			
			cp /configdb-readonly/mongod.conf "${CONFIGDIR_VOLUME}"/mongod.conf
			cp /keydir-readonly/key.txt "${CONFIGDIR_VOLUME}"/
			cp /ca-readonly/tls.key "${CONFIGDIR_VOLUME}"/tls.key
			cp /ca-readonly/tls.crt "${CONFIGDIR_VOLUME}"/tls.crt
			
			chmod 600 "${CONFIGDIR_VOLUME}"/key.txt
			chown -R 999:999 /work-dir
			chown -R 999:999 /data
			
			# Root file system is readonly but still need write and execute access to tmp
			chmod -R 777 /tmp
			`,
		},
	}

	// Set CommonServiceConfig instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, icpMongodbInstallConfigMap, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("creating icp mongodb install config map")

	if err = r.client.Create(context.TODO(), icpMongodbInstallConfigMap); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	}

	keyfileSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    metadatalabel,
			Name:      "icp-mongodb-keyfile",
			Namespace: instance.GetNamespace(),
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"key.txt": []byte("aWNwdGVzdA=="),
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
