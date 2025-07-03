### Create and Install NooBaa in `dc1` and `dc2` cluster
```bash
export KUBECONFIG=~/KUBECONFIG/noobaa-<current-cluster>.yaml
kubectl create ns noobaa
kubectl config set-context --current --namespace noobaa
noobaa install --db-volume-size-gb=10 --manual-default-backingstore=true
```
---
###  Setup In DC1
```bash
export KUBECONFIG=~/KUBECONFIG/<noobaa-dc1>.yaml
kubectl apply -f ./dc1/

# Apply below command to get the Bucket Name:
export EXTERNAL_IP=$(kubectl get svc s3 -n noobaa -o json | jq -r '.status.loadBalancer.ingress[0].ip')
export NOOBAA_ACCESS_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d')
export NOOBAA_SECRET_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d')
alias s3='AWS_ACCESS_KEY_ID=$NOOBAA_ACCESS_KEY AWS_SECRET_ACCESS_KEY=$NOOBAA_SECRET_KEY aws --endpoint https://$EXTERNAL_IP:443 --no-verify-ssl s3'
export AWS_REQUEST_CHECKSUM_CALCULATION=when_required
export AWS_RESPONSE_CHECKSUM_CALCULATION=when_required
s3 ls

# Get ehe Ip of LoadBalancer
kubectl get svc s3 -n noobaa -o json | jq -r '.status.loadBalancer.ingress[0].ip'

# Get the Secret Key and Access Key
kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d'
kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d'
```
---
###  Setup In DC2
```bash
export KUBECONFIG=~/KUBECONFIG/<noobaa-dc1>.yaml
kubectl apply -f ./dc2/

# Apply below command to get the Bucket Name:
export EXTERNAL_IP=$(kubectl get svc s3 -n noobaa -o json | jq -r '.status.loadBalancer.ingress[0].ip')
export NOOBAA_ACCESS_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d')
export NOOBAA_SECRET_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d')
alias s3='AWS_ACCESS_KEY_ID=$NOOBAA_ACCESS_KEY AWS_SECRET_ACCESS_KEY=$NOOBAA_SECRET_KEY aws --endpoint https://$EXTERNAL_IP:443 --no-verify-ssl s3'
export AWS_REQUEST_CHECKSUM_CALCULATION=when_required
export AWS_RESPONSE_CHECKSUM_CALCULATION=when_required
s3 ls

# Get ehe Ip of LoadBalancer
kubectl get svc s3 -n noobaa -o json | jq -r '.status.loadBalancer.ingress[0].ip'

# Get the Secret Key and Access Key
kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d'
kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d'
```

---
### Setup in Primary DC

```bash
## Apply Secret for dc1 with the secret key and access key of dc1
kubectl create secret generic noobaa-aws-s3-secret-dc1 \
    --from-literal=AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID> \
    --from-literal=AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY> -n noobaa
    
## Apply Secret for dc2 with the secret key and access key of dc2
kubectl create secret generic noobaa-aws-s3-secret-dc2 \
    --from-literal=AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID> \
    --from-literal=AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY> -n noobaa
    
# Now in Primary Edit the EndPoint_IP and Bucket_Name in BackingStore CRD for each dc1 and dc2


kubectl apply -f ./primary/
```


**If primary is in Kind Cluster (Need to forward port as there is no public_ip)** 

```bash
kubectl port-forward -n noobaa service/s3 10443:443

export NOOBAA_ACCESS_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d')
export NOOBAA_SECRET_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d')
alias s3='AWS_ACCESS_KEY_ID=$NOOBAA_ACCESS_KEY AWS_SECRET_ACCESS_KEY=$NOOBAA_SECRET_KEY aws --endpoint https://localhost:10443 --no-verify-ssl s3'
export AWS_REQUEST_CHECKSUM_CALCULATION=when_required
export AWS_RESPONSE_CHECKSUM_CALCULATION=when_required

s3 ls 

s3 cp ./ s3://mirror-bucket-a688e72e-c24f-4cdb-bacf-63c4fdc0c6e2/anisur2/ --recursive
```

**Disaster Scenario:**
- Insert some data to the bucket in `mirror_bucket` in primary cluster.
- Remove the `dc2` from bucket-class in primary cluster.
- Now check data still can be accessed. 

---

## Findings:
- BackingStore CRD does not have nodeName field. So, it's not possible to create POD `node/region` wise.
- So, how could we manage geolocation of data in NooBaa?
- Lack of provide nodeaffinity in BackingStore CRD.