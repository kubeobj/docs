
**Install NooBaa**
```bash
noobaa install --db-volume-size-gb=10 --disable-load-balancer=true
```

** CREATE TLS CERTIFICATE**
```bash
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:4096 \
        -keyout tls.key -out tls.crt \
        -config san.conf -extensions v3_req
```

** CREATE k8s SECRET**
```bash
kubectl -n noobaa create secret generic ingress-s3-serving-cert \
  --from-file=tls.crt --from-file=tls.key
```

** Access S3 Endpoint**
```bash
export NOOBAA_ACCESS_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d')
export NOOBAA_SECRET_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d')

alias s3='AWS_ACCESS_KEY_ID=$NOOBAA_ACCESS_KEY AWS_SECRET_ACCESS_KEY=$NOOBAA_SECRET_KEY aws --endpoint https://anisur.s3 s3'
s3 ls
```


