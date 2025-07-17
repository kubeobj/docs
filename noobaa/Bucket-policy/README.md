## Create A NooBaa account
```bash
➤ noobaa account create anisur --allow_bucket_create=false
```

## Get `anisur` Account Secrets
```bash
➤ kubectl view-secret noobaa-account-anisur
```


## Now Put a Policy to allow `anisur` to access the bucket
```bash
➤ aws s3api put-bucket-policy \
              --bucket anisur \
              --policy file://policy.json \
              --endpoint-url https://localhost:10443 \
              --no-verify-ssl

```



## Now you can access the bucket with `anisur` account
```bash

```