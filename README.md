# arctiq

## Continuous Deployment to Kubernetes Engine using Jenkins.

This multibranch repository will be used to demonstrate how Jenkins can automatically pull changes from developers and build and deploy production and development environments seamlessly.

### Setup Jenkins network and cluster

```
gcloud config set compute/zone northamerica-northeast1-a
gcloud compute networks create jenkins
gcloud container clusters create jenkins-cd \
  --network jenkins --machine-type n1-standard-2 --num-nodes 2 \
  --scopes "https://www.googleapis.com/auth/projecthosting,storage-rw,cloud-platform"
```

### Verify the cluster/nodes are up

```
gcloud container clusters list
kubectl get nodes
```

### Download Helm - The Kubernetes Package Manager

```
wget https://storage.googleapis.com/kubernetes-helm/helm-v2.9.1-linux-amd64.tar.gz
tar zxfv helm-v2.9.1-linux-amd64.tar.gz
cp linux-amd64/helm .
```

### Give Jenkins permissions in the cluster

```
kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin \
        --user=$(gcloud config get-value account)
kubectl create serviceaccount tiller --namespace kube-system
kubectl create clusterrolebinding tiller-admin-binding --clusterrole=cluster-admin \
               --serviceaccount=kube-system:tiller

```

### Initialize Helm, ensuring Tiller is installed in the cluster

```
./helm init --service-account=tiller
./helm update
./helm version
```

### Install Jenkins from Helm

```
./helm install -n cd stable/jenkins -f jenkins/values.yaml --version 0.16.6 --wait
```

### Wait for all pods to be propagated

```
kubectl get pods
```

### Print the Jenkins admin password

```
printf $(kubectl get secret cd-jenkins -o jsonpath="{.data.jenkins-admin-password}" | base64 --decode);echo
```

### Forward port 8080 to access the Jenkins Web UI

```
export POD_NAME=$(kubectl get pods -l "component=cd-jenkins-master" -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward $POD_NAME 8080:8080 >> /dev/null &
```

### Jenkins configuration

Add the credentials created above so that Jenkins may have access to the Kubernetes cluster.

Create a new item, a Multibranch Pipeline, using the HTTPS url of this GitHub repository.
