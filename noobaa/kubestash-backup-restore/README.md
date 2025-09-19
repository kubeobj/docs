### NooBaa Backup/Restore with kubeStash

1. **Take Backup of all secrets in which're related to Noobaa**
2. **In a new cluster**
   1. Install noobaa operator only
   2. Restore all secrets
   3. Restore the backuped PG data to the new or old targeted database
   4. Apply the Noobaa CRD
   5. Restore all crds related to NooBaa
   6. Now See the **Magic**
   

### Workaround for Noobaa Backup/Restore

```bash
kubectl create secret generic noobaa-external-pg-db \
              --namespace=noobaa \
              --from-literal=db_url='postgres://postgres:Lg73qGD.YMVB838F@10.2.0.81:5432/nbcore'
```

```bash 
# IN cluster where the postgres is running
kubectl exec -it pod/nbcore-postgres-0 -- pg_dump nbcore -f /tmp/test.db -F custom
kubectl cp nbcore-postgres-0:/tmp/test.db ./mcg_new.bck


kubectl exec -it nbcore-postgres-0 -- pg_restore -d nbcore /tmp/test.db -c

kubectl exec -it nbcore-postgres-0 -- \
  psql -d nbcore -c "DROP TABLE IF EXISTS public.nodes;"
```
