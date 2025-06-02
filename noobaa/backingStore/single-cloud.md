### Install noobaa-operator-cli

https://github.com/noobaa/noobaa-operator/releases/tag/v5.18.1

### Install to Kubernetes

```bash
kubectl create ns noobaa
kubectl config set-context --current --namespace noobaa
noobaa install
```

### Default BackingStores check (type: pv-pool)

```bash
kubectl describe noobaa


kubectl port-forward -n noobaa service/s3 10443:443

export NOOBAA_ACCESS_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d')
export NOOBAA_SECRET_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d')
alias s3='AWS_ACCESS_KEY_ID=$NOOBAA_ACCESS_KEY AWS_SECRET_ACCESS_KEY=$NOOBAA_SECRET_KEY aws --endpoint https://localhost:10443 --no-verify-ssl s3'
export AWS_REQUEST_CHECKSUM_CALCULATION=when_required
export AWS_RESPONSE_CHECKSUM_CALCULATION=when_required

s3 ls 

s3 cp example.txt s3://<bucket_name>/

s3 ls
```

### Create a new AWS S3 `BackingStores`

```bash
noobaa backingstore create aws-s3 aws-bs --region=us-east-1 --secret-name aws-s3-secret --target-bucket noobaa
noobaa bucketclass create placement-bucketclass aws-bc --backingstores=aws-bs
noobaa obc create aws-claim --bucketclass=aws-bc

s3 ls
```



