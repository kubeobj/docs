### We Can increase total Storage VOLUME 
```yaml
  pvPool:
    numVolumes: 1 # Increasing num of volumes
```

### See any bucket Storage Availability status:
```bash
$ noobaa bucket status anisur-cf9cd250-1376-42a8-b783-ee29bd674fc0

Bucket status:
  Bucket                 : anisur-cf9cd250-1376-42a8-b783-ee29bd674fc0
  OBC Namespace          : noobaa
  OBC BucketClass        : sample-bucket-class
  Type                   : REGULAR
  Mode                   : OPTIMAL
  ResiliencyStatus       : OPTIMAL
  QuotaStatus            : QUOTA_NOT_SET
  Num Objects            : 84933
  Data Size              : 3.289 GB
  Data Size Reduced      : 206.638 MB
  Data Space Avail       : 15.189 GB
  Num Objects Avail      : 9007199254740991

```

# After Increasing numVolumes

```bash
$ noobaa bucket status anisur-cf9cd250-1376-42a8-b783-ee29bd674fc0

Bucket status:
  Bucket                 : anisur-cf9cd250-1376-42a8-b783-ee29bd674fc0
  OBC Namespace          : noobaa
  OBC BucketClass        : sample-bucket-class
  Type                   : REGULAR
  Mode                   : OPTIMAL
  ResiliencyStatus       : OPTIMAL
  QuotaStatus            : QUOTA_NOT_SET
  Num Objects            : 86682
  Data Size              : 3.459 GB
  Data Size Reduced      : 206.642 MB
  Data Space Avail       : 30.640 GB
  Num Objects Avail      : 9007199254740991

```

