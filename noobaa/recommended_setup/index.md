## HA NooBaa Setup

**Requirements:**
- **Cluster Size:** Use `3–5` Kubernetes worker nodes. 
  - `Three` is the minimum for quorum and tolerating a single node failure. 
  - `Five` is preferred for better redundancy and rebuild speed.
- **Storage Class**: Use a distributed, replicated storage backend. 
  - **Recommended:**
    - longhorn (node-level replication) 
    - rook-ceph (replication + erasure coding)
  - **Replication vs. Erasure Coding:**
    - **Replication:** 3× copies; 10 TB usable needs ~30 TB raw.
    - **Erasure Coding:** ~1.5× overhead; 10 TB usable needs ~15 TB raw.

For large data, erasure coding is more space-efficient.

**How it'll ensure HA:**
- **Data Plane** (Main Storage):
  - Each PV used by NooBaa will be backed by a storage system that replicates or erasure-codes data across multiple nodes. 
  - If one node fails, replicas/parity blocks on other nodes keep data available without manual intervention.
- **Control Plane** (NooBaa Brain):
 - The Brain (metadata) is the only single point of failure in NooBaa’s architecture. 
 - Replace NooBaa’s default single Postgres pod with KubeDB by AppsCode Managed Postgres in HA mode. 
 - KubeStash by AppsCode PITR (Point-In-Time Recovery) backup/restore feature.

**S3 Endpoint Scaling:**
- The S3 endpoint pods are stateless and can be scaled `horizontally` to handle traffic from `millions of clients`.
- Place them behind a Kubernetes `Ingress or LoadBalancer` with `TLS` termination for secure access.

- As for each PV data will be replicated across multiple nodes, so if one node goes down, the data will still be available on other nodes.
- We'll use KubeDB Managed Postgres for Brain, this is the only single point of failure in NooBaa, so using KubeDB Managed Postgres will ensure that the Brain is highly available, we can take PITR backup easily.
- While the s3 endpoint client increase day by day, we can scale the NooBaa endpoint deployment horizontally to handle the increased load.

**Cross-Cluster Disaster Recovery:** (Optional)
- For extreme resilience (geo-redundancy), deploy two or more NooBaa clusters in different regions or datacenters.
- Use NooBaa Bucket Mirroring to replicate buckets between clusters. 
  - If one cluster is down, the mirrored bucket in the second cluster remains accessible.


## Setup with HA, Ingress, TLS, and Mirroring

![design.png](images/design.png)

**Install KubeDB**

**Install KubeStash**

**Create a HA postgres With PITR enable**

```bash
kubectl apply -f .yamls/ha-postgres.yaml
```

**Create nbcore DB**
```bash
kubectl exec -it nbcore-postgres-0 -- psql -U postgres -d postgres -c "CREATE DATABASE nbcore WITH LC_COLLATE = 'C' TEMPLATE template0;"
````


**Prepare DB URL**
```bash
kubectl view-secret nbcore-postgres-auth
password=''
username='postgres'

postgres://postgres:password@nbcore-postgres.noobaa.svc.cluster.local:5432/nbcore
```

**CREATE TLS CERTIFICATE**
```bash
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:4096 \
        -keyout tls.key -out tls.crt \
        -config san.conf -extensions v3_req
```

**CREATE TLS SECRET**
```bash
kubectl -n noobaa create secret generic noobaa-s3-serving-cert \
  --from-file=tls.crt --from-file=tls.key

kubectl -n noobaa create secret generic noobaa-mgmt-serving-cert \
  --from-file=tls.crt --from-file=tls.key
```



**Install NooBaa Operator**
Download your desired file, extract the binary and place it in your PATH. https://github.com/noobaa/noobaa-operator/releases/tag/v5.18.1
```bash
noobaa install --disable-load-balancer=true  --postgres-url='<db-url>' --manual-default-backingstore=true
```

**Install Default BackingStore**
```bash
apiVersion: noobaa.io/v1alpha1
kind: BackingStore
metadata:
  name: noobaa-default-backing-store
  namespace: noobaa
spec:
  pvPool:
    numVolumes: 1
    resources:
      requests:
        storage: 16Gi
        memory: 2Gi
        cpu: 1000m
      limits:
        memory: 2Gi
        cpu: 1000m
    secret: {}
  type: pv-pool
```

> Note: It'll use when you create a bucket using noobaa cli.

**Enable Ingress**
```bash
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
````

***Create TLS Secret for ingress***
```bash
kubectl -n noobaa create secret generic ingress-s3-serving-cert \
  --from-file=tls.crt --from-file=tls.key
````

**Create Ingress Resource**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: noobaa-ingress
  namespace: noobaa
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - anisur.s3          # must match spec.rules[*].host
      secretName: ingress-s3-serving-cert           # the Secret we just created
  rules:
    - host: anisur.s3
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: s3
                port:
                  number: 80
```

**Access S3 Endpoints**
```bash
export NOOBAA_ACCESS_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d')
export NOOBAA_SECRET_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d')
alias s3='AWS_ACCESS_KEY_ID=$NOOBAA_ACCESS_KEY AWS_SECRET_ACCESS_KEY=$NOOBAA_SECRET_KEY aws --endpoint https://anisur.s3 s3'
export AWS_REQUEST_CHECKSUM_CALCULATION=when_required
export AWS_RESPONSE_CHECKSUM_CALCULATION=when_required

aws configure set ca_bundle tls.crt
s3 ls 
```
---

### Test Mirroring Feature of NooBaa
**Scale the s3 endpoint deployment(HA)**
```bash
kubectl scale deployment noobaa-endpoint --replicas=3
```

**Create BackingStore for Mirroring**
```yaml
apiVersion: noobaa.io/v1alpha1
kind: BackingStore
metadata:
  name: local-pv-pool-1
  namespace: noobaa
spec:
  pvPool:
    numVolumes: 1
    resources:
      requests:
        storage: 16Gi
        memory: 2Gi
        cpu: 1000m
      limits:
        memory: 2Gi
        cpu: 1000m
  type: pv-pool
---
apiVersion: noobaa.io/v1alpha1
kind: BackingStore
metadata:
  name: local-pv-pool-2
  namespace: noobaa
spec:
  pvPool:
    numVolumes: 1
    resources:
      requests:
        storage: 16Gi
        memory: 2Gi
        cpu: 1000m
      limits:
        memory: 2Gi
        cpu: 1000m
  type: pv-pool
```

Here, `local-pv-pool-1` and `local-pv-pool-2` are the backing stores that will be used for mirroring. So, If one get down the other will served the data.

**Create BucketClass for Mirroring**
```yaml
apiVersion: noobaa.io/v1alpha1
kind: BucketClass
metadata:
  name: mirror-bucket-class
spec:
  placementPolicy:
    tiers:
      - backingStores:
          - local-pv-pool-1
          - local-pv-pool-1
        placement: Mirror
```

** Now creat a OBC with mirror bucketclass**:
```bash
apiVersion: objectbucket.io/v1alpha1
kind: ObjectBucketClaim
metadata:
  name: mirror-obc
spec:
  generateBucketName: mirror-bucket
  storageClassName: noobaa.noobaa.io
  additionalConfig:
    bucketclass: mirror-bucket-class
```

Here, `mirror-bucket` is the name of the bucket that will be created. And we can claiming bucket as like PVC as kubernetes native way.  
So, It created a bucket with mirroring policy.
```bash
➤ s3 ls
2025-08-14 15:57:53 mirror-bucket-0a18893b-5c89-494f-848f-3113f353e084
2025-08-14 15:28:45 first.bucket
```

**Upload a large to the bucket**
```bash
➤ s3 cp dune.mkv s3://mirror-bucket-0a18893b-5c89-494f-848f-3113f353e084/movies/
upload: ./dune.mkv to s3://mirror-bucket-0a18893b-5c89-494f-848f-3113f353e084/movies/dune.mkv

➤ s3 ls s3://mirror-bucket-0a18893b-5c89-494f-848f-3113f353e084/ --recursive
2025-08-14 16:12:04 3454831563 movies/dune.mkv

```

**Now, If one of the backing store is down, the other will serve the data.**

I've uploaded 12 GB file to the bucket. Now Going to delete one of the backing store while a process is downloading the file.

```bash
➤ s3 cp s3://mirror-bucket-0a18893b-5c89-494f-848f-3113f353e084/movies/ ./noobaa-download/

## From another terminal, delete one of the backing store, I am going to delete local-pv-pool-dc2

➤ kubectl delete backingstore.noobaa.io/local-pv-pool-2
```

**Increase the Storage of a backing store**
```bash
apiVersion: noobaa.io/v1alpha1
kind: BackingStore
spec:
  pvPool:
    numVolumes: 1 ## Simply, increase the numVolumes to increase the storage

```
> Note: Once you increase the storage, you can't descre the backing store. In that case, you've to delete the backing store.


**Testing S3 Bucket Policy**

```bash
➤ noobaa account create anisur --allow_bucket_create=false
➤ kubectl view-secret noobaa-account-anisur -a
➤ export NOOBAA_ACCESS_KEY=
➤ export NOOBAA_SECRET_KEY=

➤ s3 ls s3://mirror-bucket-0a18893b-5c89-494f-848f-3113f353e084/ --recursive

An error occurred (AccessDenied) when calling the ListObjectsV2 operation: Access Denied
# No access of new account 

# Now, Create a bucket policy to allow the new account to access the bucket
# Below operation need to be done in the different terminal using root account access
➤ aws s3api put-bucket-policy \
              --bucket anisur \
              --policy file://policy.json \
              --endpoint-url --endpoint https://anisur.s3
              
➤ s3 ls s3://mirror-bucket-0a18893b-5c89-494f-848f-3113f353e084/ --recursive
2025-08-14 16:27:10 3454831563 movies/another_dune.mkv
2025-08-14 16:12:04 3454831563 movies/dune.mkv
2025-08-14 16:28:46 3454831563 movies/dune2.mkv
2025-08-14 16:29:57 3454831563 movies/dune3.mkv
```


**Example deploying a HTML page**

Here, I'm going to deploy a static HTML page to the bucket.
```bash
➤ aws s3 mb s3://web --endpoint-url https://anisur.s3
make_bucket: web

➤ aws s3 cp index.html s3://web/index.html \
        --content-type 'text/html' --endpoint-url https://anisur.s3
upload: ./index.html to s3://web/index.html  

➤ aws s3api put-bucket-policy \
                    --bucket web \
                    --policy file://web-policy.json \
                    --endpoint-url https://anisur.s3

➤ aws s3api put-bucket-website --bucket web --website-configuration '{
      "IndexDocument": { "Suffix": "index.html" }
    }' --endpoint-url https://anisur.s3

# Now, visit https://anisur.s3/web
```
