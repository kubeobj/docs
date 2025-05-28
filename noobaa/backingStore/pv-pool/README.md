
## BucketClass 
FIELDS:
namespacePolicy	<Object>
NamespacePolicy specifies the namespace policy for the bucket class

placementPolicy	<Object>
PlacementPolicy specifies the placement policy for the bucket class

quota	<Object>
Quota specifies the quota configuration for the bucket class

replicationPolicy	<string>
ReplicationPolicy specifies a json of replication rules for the bucketclass

### `bucketclasses.noobaa.io.spec.placementPolicy`

In NooBaa, a `BucketClass` defines how data is stored and managed across different storage backends. The `spec.placementPolicy` field within a `BucketClass` specifies the strategy for placing data across one or more backing stores.

#### 1. `tiers`
This is a list that defines one or more tiers for data placement. Each tier includes:

* **`backingStores`**: An array of backing store names where data will be stored.
* **`placement`**: The strategy used to distribute data across the specified backing stores.

#### 2. `placement` Strategies
* **`Spread`**: Distributes data evenly across all specified backing stores. This approach balances storage utilization and can enhance performance.
```yaml
  placementPolicy:
    tiers:
      - backingStores:
          - store1
          - store2
        placement: Spread
```
* **`Mirror`**: Replicates data to all specified backing stores. This strategy provides high availability and redundancy.([Red Hat Documentation][5])
```yaml
  placementPolicy:
    tiers:
      - backingStores:
          - store1
          - store2
        placement: Mirror
```
* **`Single`**: Stores all data in a single backing store. This is a straightforward approach without redundancy.
```yaml
  placementPolicy:
    tiers:
      - backingStores:
          - store1
        placement: Single
```
