### Introduction
`SeaweedFS` is a simple and highly scalable distributed file system. There are two objectives:
- to store billions of files!
- to serve the files fast

Instead of managing all file metadata in a central master, the central master only manages volumes on volume servers, and these volume servers manage files and their metadata.

### How it works:
- Weed server `-s3` will start a `master`, a `volume server`, a `filer`, and the `S3 gateway`.
- The `master` server is responsible for assigning file ids which volume to store the objects and managing the volume servers and their volumes. Unlike Other's master server only for volume manages not individual files.
- The `volume server` is responsible for storing the files and their metadata. (sends back heartbeats to the lead Master service)
- The `filer` server is responsible to provide file system-like interface with support for directories, POSIX attributes, and APIs `S3`.
- The `S3 gateway` is responsible for providing an S3-compatible API to the `SeaweedFS` file system.

### Feature Differences With Real S3
- https://github.com/seaweedfs/seaweedfs/wiki/Amazon-S3-API#feature-difference
- It Supports `fake` directories, while `S3` not.
- It not allow more than `/` as a delimiter, while `S3` allow more than one.



### Seaweed Admin 
- `Closed-source` mainly they do feautres like backup/recovery and a management console.
- One-Click Restore
- https://seaweedfs.com/docs/admin/features/

Basically, we can do that stuff using `kubestash`.


### RUN in locally

```bash
$ weed server -dir=./local/data \
                  -master.volumeSizeLimitMB=1024 \
                  -master.volumePreallocate=false \
                  -s3 \
                  -s3.port=9000 \
                  -s3.config=./local/s3_config.jso
```

**Create Bucket using weed shell**
```bash
$ echo "s3.bucket.create -name anisur" | weed shell
```

**AWS CLI**
```bash
$ docker run --rm -it -v ./aws:/root/.aws amazon/aws-cli configure
$ docker run --rm -it -v ./aws:/root/.aws amazon/aws-cli configure set default.s3.signature_version s3v4
$ docker run --rm -it --network=host -v ./aws:/root/.aws amazon/aws-cli --endpoint-url http://localhost:9000 s3 ls
```



### References:
- https://blog.jklug.work/posts/seaweedfs/
- https://github.com/JBris/seaweedfs-s3-trial
- 
