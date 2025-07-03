### NamespaceStore
A NamespaceStore (used in Namespace Buckets) is a pointer to an external object storage.
NooBaa does not store data itself here — instead, it proxies requests to the external storage.

**Create S3 BackingStore**
- Create Secret
```bash
kubectl create secret generic aws-s3-secret \
        --from-literal=AWS_ACCESS_KEY_ID=<> \
        --from-literal=AWS_SECRET_ACCESS_KEY=<> \
        -n demo
```

- Create `NamespaceStore` CR

```yaml
apiVersion: noobaa.io/v1alpha1
kind: NamespaceStore
metadata:
  name: s3-namespace-bs
spec:
  awsS3:
    region: us-east-1
    secret:
      name: aws-s3-secret
      namespace: noobaa
    targetBucket: namespace-noobaa
  type: aws-s3
```

**Create AZURE BackingStore**

- Create Secret
```bash
kubectl create secret generic azure-secret \
  --from-literal=AccountName=kubstashqa \
  --from-literal=AccountKey=<Account_Key>
```

- Create `NamespaceStore` CR

```yaml
apiVersion: noobaa.io/v1alpha1
kind: NamespaceStore
metadata:
  name: gcs-namespace-bs
spec:
  googleCloudStorage:
    secret:
      name: gcs-secret
      namespace: noobaa
    targetBucket: namespace-noobaa
  type: google-cloud-storage
```

**namespacePolicy**

```yaml
apiVersion: noobaa.io/v1alpha1
kind: BucketClass
metadata:
  name: multi-namespace-bc
spec:
  namespacePolicy:
    type: Multi
    multi:
      readResources:
        - azure-namespace-bs
      writeResource: s3-namespace-bs
```

Here,
- There are no options for replication. 
- ONly we can do read from many and write only one namespace store.
