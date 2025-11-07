---
title: Readme
menu:
  docs_v2025.10.31:
    identifier: readme-concepts
    name: Readme
    parent: concepts
    weight: 1
menu_name: docs_v2025.10.31
section_menu_id: concepts
url: /docs/v2025.10.31/concepts/
aliases:
- /docs/v2025.10.31/concepts/README/
info:
  cli: v0.20.0
  installer: v2025.10.31
  version: v2025.10.31
---

# Concepts
Concepts help you to learn about the different parts of the NooBaa and the abstractions it uses.

This concept section is divided into the following modules:

## What is NooBaa?
- [Overview](/docs/v2025.10.31/concepts/what-is-noobaa/overview/) provides an introduction to NooBaa. It also gives an overview of the features it provides.

## Declarative API

- [NooBaa](/docs/v2025.10.31/concepts/crds/noobaa/): This is the central CRD for a NooBaa deployment. It represents the entire NooBaa system and orchestrates all its components. Administrators interact with this resource to install, configure, and manage the NooBaa platform within a Kubernetes cluster.
- [BackingStore](/docs/v2025.10.31/concepts/crds/backingstore/): A `BackingStore` represents underlying storage resources that NooBaa can use to store data. These can be cloud-based storage services (like AWS S3, Azure Blob Storage, or Google Cloud Storage) or on-premises storage systems (like a PV pool). `BackingStore`s are used to create buckets and provide the physical storage layer.
- [NamespaceStore](/docs/v2025.10.31/concepts/crds/namespacestore/): A `NamespaceStore` is a special type of `BackingStore` that represents an entire external bucket or a prefix within a bucket. Instead of using the external storage as a backing for NooBaa's own data services (like deduplication and compression), a `NamespaceStore` provides a direct, transparent connection to the external resource. This is useful for data federation and accessing existing data in place.
- [BucketClass](/docs/v2025.10.31/concepts/crds/bucketclass/): A `BucketClass` defines policies for how data is placed and managed across one or more `BackingStore`s. It can specify data placement policies like mirroring or spreading data across different stores. When a user requests a new bucket, they can specify a `BucketClass` to control the behavior and characteristics of their bucket.
- [ObjectBucketClaim](/docs/v2025.10.31/concepts/crds/objectbucketclaim/): An `ObjectBucketClaim` (OBC) is a Kubernetes-native way for a user or application to request a new object storage bucket, similar to how a `PersistentVolumeClaim` (PVC) requests block or file storage. When an OBC is created, the NooBaa operator provisions a corresponding `ObjectBucket` and makes the connection details available to the application via a `Secret` and a `ConfigMap`.
