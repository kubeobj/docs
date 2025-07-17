## External TLS Certificate Generation
** CREATE TLS CERTIFICATE**
```bash
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:4096 \
        -keyout tls.key -out tls.crt \
        -config san.conf -extensions v3_req
```

** CREATE k8s SECRET**
```bash
kubectl -n noobaa create secret generic noobaa-s3-serving-cert \
  --from-file=tls.crt --from-file=tls.key
```

** Check WITH TLS**
```bash
aws s3 configure
aws s3 ls --endpoint https://anisur.s3:<s3-service-port>
aws s3 mb s3://web
aws s3 cp index.html s3://web/index.html \
  --content-type 'text/html'


aws s3api put-bucket-policy \
              --bucket web \
              --policy file://policy.json \
              --endpoint-url https://anisur.s3
  
aws s3api put-bucket-website --bucket web --website-configuration '{
    "IndexDocument": { "Suffix": "index.html" }
  }' --endpoint-url https://anisur.s3:<s3-service-port>



https://anisur.s3:<s3-service-port>/web/
```

