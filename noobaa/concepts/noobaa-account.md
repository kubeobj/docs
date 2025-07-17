### NooBaaAccount

```bash
noobaa -n noobaa account create account1 --allow_bucket_create --default_resource bs1
```
It'll create a new account named `account1` with the ability to create buckets and and bucket will store data in backing store `bs1`.

Notes:
- Here, We can't use bucketclass, that means, there is no mirroing or spreading of data.


