### BackingStores
A BackingStore represents physical storage where data is stored.



### BucketClass
Bucket class is a CRD representing a class of buckets that defines tiering policies 
and data placements for an Object Bucket Class (OBC).

**BucketClass type**
- Standard: data will be consumed by a Multicloud Object Gateway (MCG), deduped, compressed and encrypted.
- Namespace: data is stored on the NamespaceStores without performing de-duplication, compression or encryption.


**Placement Policy**
- Spread: allows spreading of the data across the chosen resources.
- Mirror: allows full duplication of the data across the chosen resources.

### 





### NamespaceStore

A NamespaceStore (used in Namespace Buckets) is a pointer to an external object storage. 
NooBaa does not store data itself here — instead, it proxies requests to the external storage.