### Install Noobaa for Kuebrntess
```bash
# Prepare namespace and set as current (optional)
kubectl create ns noobaa
kubectl config set-context --current --namespace noobaa

# Install the operator and system on your cluster:
noobaa install

# You can always get system status and information with:
noobaa status
```
### Get S3 access keys and setup s3 CLI as client

```bash
export NOOBAA_ACCESS_KEY=$(kubectl get secret noobaa-admin  -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d')
export NOOBAA_SECRET_KEY=$(kubectl get secret noobaa-admin  -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d')
alias s3='AWS_ACCESS_KEY_ID=$NOOBAA_ACCESS_KEY AWS_SECRET_ACCESS_KEY=$NOOBAA_SECRET_KEY aws --endpoint https://localhost:10443 --no-verify-ssl s3'
export AWS_REQUEST_CHECKSUM_CALCULATION=when_required
export AWS_RESPONSE_CHECKSUM_CALCULATION=when_required
```

### Copy file to generated bucket
```bash
s3 ls s3://
s3 cp ./ s3://<bucket-name>/ --recursive
```

https://chatgpt.com/share/6821d2ba-3cd0-8001-bc87-79dd85848ad5