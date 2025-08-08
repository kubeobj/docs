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
**Install Longhorn** 
- Enable volumeSnapshot feature to cluster. [Link](https://github.com/kubernetes-csi/external-snapshotter?tab=readme-ov-file)
```bash
kubectl kustomize https://github.com/kubernetes-csi/external-snapshotter/client/config/crd | kubectl create -f -
kubectl -n kube-system kustomize https://github.com/kubernetes-csi/external-snapshotter/deploy/kubernetes/snapshot-controller | kubectl create -f -
kubectl kustomize https://github.com/kubernetes-csi/external-snapshotter/deploy/kubernetes/csi-snapshotter | kubectl create -f -
```
- Set cloud information from longhorn UI.
**Install KubeDB**
**Install KubeStash**

**Create a HA postgres**
