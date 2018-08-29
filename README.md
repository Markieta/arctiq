# arctiq

## Continuous Deployment to Kubernetes Engine using Jenkins.

This multibranch repository will be used to demonstrate how Jenkins can automatically pull changes from developers and build and deploy production and development environments seamlessly.

```
gcloud config set compute/zone northamerica-northeast1-a
gcloud compute networks create jenkins
gcloud container clusters create jenkins-cd \
  --network jenkins --machine-type n1-standard-2 --num-nodes 2 \
  --scopes "https://www.googleapis.com/auth/projecthosting,storage-rw,cloud-platform"
```

```
gcloud container clusters list
kubectl get nodes
```

```
wget https://storage.googleapis.com/kubernetes-helm/helm-v2.9.1-linux-amd64.tar.gz
tar zxfv helm-v2.9.1-linux-amd64.tar.gz
cp linux-amd64/helm .
```
```
kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin \
        --user=$(gcloud config get-value account)
kubectl create serviceaccount tiller --namespace kube-system
kubectl create clusterrolebinding tiller-admin-binding --clusterrole=cluster-admin \
               --serviceaccount=kube-system:tiller

```

```
./helm init --service-account=tiller
./helm update
./helm version
```

```
./helm install -n cd stable/jenkins -f jenkins/values.yaml --version 0.16.6 --wait
```
