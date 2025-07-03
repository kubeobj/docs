### Install NOOBAA
```bash
export KUBECONFIG=~/KUBECONFIG/noobaa-<current-cluster>.yaml
kubectl create ns noobaa
kubectl config set-context --current --namespace noobaa
noobaa install --db-volume-size-gb=10 --manual-default-backingstore=true
```

### Access NooBaa by using External LoadBalancer IP=`<EXTERNAL_IP>`
```bash
export EXTERNAL_IP=$(kubectl get svc s3 -n noobaa -o json | jq -r '.status.loadBalancer.ingress[0].ip')
export NOOBAA_ACCESS_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d')
export NOOBAA_SECRET_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d')
alias s3='AWS_ACCESS_KEY_ID=$NOOBAA_ACCESS_KEY AWS_SECRET_ACCESS_KEY=$NOOBAA_SECRET_KEY aws --endpoint https://$EXTERNAL_IP:443 --no-verify-ssl s3'
export AWS_REQUEST_CHECKSUM_CALCULATION=when_required
export AWS_RESPONSE_CHECKSUM_CALCULATION=when_required
s3 ls
```