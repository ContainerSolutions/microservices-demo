# ServiceBroker for Redis

## Service Catalog

Install [service
catalog](https://kubernetes.io/docs/concepts/extend-kubernetes/service-catalog/).
The following instructions use helm to do so. More information can be found
[here](https://kubernetes.io/docs/tasks/service-catalog/install-service-catalog-using-helm/)

Configure Tiller to have `cluster-admin access`

```shell
kubectl create clusterrolebinding tiller-cluster-admin \
    --clusterrole=cluster-admin \
    --serviceaccount=kube-system:default
```

Add service catalog helm repository and install service catalog:

```shell
helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
helm install svc-cat/catalog --name catalog --namespace catalog
```

## Open Service Broker for Google Cloud Platform

We're going to use helm to install the [Service Broker](https://github.com/GoogleCloudPlatform/gcp-service-broker/tree/feature/helm/deployments/helm/gcp-service-broker). The helm chart for the Service Broker is not yet merged to master.

Get the helm chart:

```shell
git clone git@github.com:GoogleCloudPlatform/gcp-service-broker.git
cd gcp-service-broker
git checkout develop
```

Create a Service Account for the Service Broker as shown [here](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/installation.md#create-a-root-service-account).

Create a key for the Service Account and download it as JSON and save it as
`osservicebroker-key.json`

Use the following script to install the Service Broker:

```bash
#!/bin/bash

CHART_NAME=gcp-service-broker
CHART_DIR=./deployments/helm/gcp-service-broker
SERVICE_ACCOUNT_KEY=./osservicebroker-key.json
BROKER_VERSION=4.2.2

helm package \
    --app-version=$BROKER_VERSION \
    --dependency-update \
    --version=$BROKER_VERSION \
    $CHART_DIR

helm upgrade -i gcp-service-broker \
    --set-file broker.service_account_json=./osservicebroker-key.json \
    --set broker.password=Ohdei3 \
   ${CHART_NAME}-${BROKER_VERSION}.tgz
```

### Verify

Pods are running:

```shell
% kubectl get pods -n catalog
NAME                                                 READY     STATUS    RESTARTS   AGE
catalog-catalog-apiserver-58cdffb8c5-49c7x           2/2       Running   0          44m
catalog-catalog-controller-manager-bcc4879cf-jkdd6   1/1       Running   0          44m
gcp-service-broker-6b96db74dc-wn6xw                  1/1       Running   3          2m
gcp-service-broker-mysql-59b5bcbdc4-bhwjq            1/1       Running   0          2m
```

The broker is registered to the service catalog:

```shell
% svcat get brokers
         NAME                                  URL                           STATUS
+--------------------+-----------------------------------------------------+--------+
  gcp-service-broker   http://gcp-service-broker.catalog.svc.cluster.local   Ready
```

## Create an Instance

```yaml
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ServiceInstance
metadata:
  name: mysql01
spec:
  clusterServiceClassExternalName: google-cloudsql-mysql
  clusterServicePlanExternalName: mysql-db-f1-micro
```

And then `kubectl apply -f` the file above. You should get:

```shell
% gcloud sql instances list
NAME                       DATABASE_VERSION  LOCATION        TIER              PRIMARY_ADDRESS  PRIVATE_ADDRESS  STATUS
gsb-1-1552477537279273624  MYSQL_5_7         us-central1-b   db-f1-micro       104.154.121.178  -                PENDING_CREATE

```

## Bind an Instance

```yaml
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ServiceBinding
metadata:
  name: mysql01-binding
spec:
  instanceRef:
    name: mysql01

```

And then `kubectl apply -f` the file above. You should get:

```shell
% kubectl get servicebindings
NAME              SERVICE-INSTANCE   SECRET-NAME       STATUS    AGE
mysql01-binding   mysql01            mysql01-binding   Ready     2m
```

The secret contains the information about the CloudSQL instance:

```shell
% kubectl describe secret mysql01-binding
Name:         mysql01-binding
Namespace:    catalog
Labels:       <none>
Annotations:  <none>

Type:  Opaque

Data
====
host:                      15 bytes
Name:                      20 bytes
Password:                  44 bytes
PrivateKeyData:            3156 bytes
Username:                  16 bytes
Email:                     69 bytes
last_master_operation_id:  0 bytes
region:                    11 bytes
Sha1Fingerprint:           40 bytes
UniqueId:                  21 bytes
database_name:             25 bytes
UriPrefix:                 0 bytes
instance_name:             25 bytes
uri:                       131 bytes
CaCert:                    1146 bytes
ClientCert:                1232 bytes
ClientKey:                 1674 bytes
ProjectId:                 24 bytes
```

For example, here's the URI of the CloudSQL instance:

```shell
 kubectl get secret mysql01-binding -o go-template --template="{{.data.uri}}" | base64 -d                                     
mysql://sb15524784435837:lQWBvX2aada7aRS_ysnEqmm4kcS_4EYpfXxBcS2vRPA%3D@104.154.121.178/gsb-2-1552477537279343171?ssl_mode=required
```

## Deploy ProductCatalogService with MySQL

TODO
