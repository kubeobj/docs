### Create and Install NooBaa in `dc1` and `dc2` cluster
```bash
export KUBECONFIG=~/KUBECONFIG/noobaa-<current-cluster>.yaml
kubectl create ns noobaa
kubectl config set-context --current --namespace noobaa
noobaa install --db-volume-size-gb=10 --manual-default-backingstore=true
```

### Access dc1,dc2 NooBaa by using External LoadBalancer IP=`<EXTERNAL_IP>`
```bash
export EXTERNAL_IP=$(kubectl get svc s3 -n noobaa -o json | jq -r '.status.loadBalancer.ingress[0].ip')
export NOOBAA_ACCESS_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d')
export NOOBAA_SECRET_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d')
alias s3='AWS_ACCESS_KEY_ID=$NOOBAA_ACCESS_KEY AWS_SECRET_ACCESS_KEY=$NOOBAA_SECRET_KEY aws --endpoint https://$EXTERNAL_IP:443 --no-verify-ssl s3'
export AWS_REQUEST_CHECKSUM_CALCULATION=when_required
export AWS_RESPONSE_CHECKSUM_CALCULATION=when_required
s3 ls
```


###  Setup In DC1

**BackingStore in DC1**

```bash
kubectl apply -f ./dc2/local-pv-pool-dc2.yaml
```


### Main setup in DC1

**BackingStore in DC2**

```bash
kubectl apply -f ./dc1/local-pv-pool-dc2.yaml
```

***Secret in DC1 with DC2's Credentials***

```bash
kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d'
kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d'

# Get DC2 S3 secrets access key and ID;
# Now create a secret in DC1 with DC2's credentials

kubectl create secret generic aws-s3-secret \ 
    --from-literal=AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID> \
    --from-literal=AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY> -n noobaa

```




## Findings:
- BackingStore CRD does not have nodeName field. So, it's not possible to create POD `node/region` wise.
- So, how could we manage geolocation of data in NooBaa?
- Lack of provide nodeaffinity in BackingStore CRD.