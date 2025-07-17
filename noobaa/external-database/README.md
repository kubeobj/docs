# External Postgresql DB support
**Notes:**
- Noobaa only support version: `postgresql 15`. (Don't know why)

** Install Noobaa Operator**
```bash
noobaa operator install

noobaa crd create
````

**Create nbcore DB**
```bash
kubectl exec -it nbcore-postgres-0 -- psql -U postgres -d postgres -c "CREATE DATABASE nbcore WITH LC_COLLATE = 'C' TEMPLATE template0;"

````

**Prepare DB URL**
```bash
➤ kubectl view-secret nbcore-postgres-auth
password=''
username='postgres'

postgres://postgres:password@nbcore-postgres.noobaa.svc.cluster.local:5432/nbcore
```

**Create Auth Secret**
```bash
kubectl create secret generic noobaa-external-pg-db \
        --namespace=noobaa \
        --from-literal=db_url='postgres://postgres:password@nbcore-postgres.noobaa.svc.cluster.local:5432/nbcore'
```

