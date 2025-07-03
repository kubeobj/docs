# External Postgresql DB support
**Notes:**
- Noobaa only support version: `postgresql 15`. (Don't know why)

**Create Auth Secret**
```bash
kubectl create secret generic noobaa-external-pg-db \
        --namespace=noobaa \
        --from-literal=db_url=''
```

