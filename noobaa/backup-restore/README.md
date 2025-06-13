
```bash
kubectl get secret noobaa-root-master-key-backend -o yaml > ./yamls/noobaa-root-master-key-backend.yaml
kubectl get secret noobaa-root-master-key-volume -o yaml >  ./yamls/noobaa-root-master-key-volume.yaml
kubectl get secret noobaa-admin -o yaml > ./yamls/noobaa-admin.yaml
kubectl get secret noobaa-db -o yaml > ./yamls/noobaa-db.yaml
kubectl get secret noobaa-operator -o yaml > ./yamls/noobaa-operator.yaml
kubectl get secret noobaa-server -o yaml > ./yamls/noobaa-server.yaml
kubectl get secret noobaa-endpoints -o yaml > ./yamls/noobaa-endpoints.yaml
kubectl get secret aws-s3-secret -o yaml > ./yamls/aws-s3-secret.yaml

kubectl exec -it noobaa-db-pg-0 -- pg_dump nbcore -f /tmp/test.db -F custom
kubectl cp noobaa-db-pg-0:/tmp/test.db ./mcg.bck
```


```bash
kubectl scale deployment noobaa-operator  --replicas=0
kubectl scale deployment noobaa-endpoint  --replicas=0
kubectl scale sts noobaa-core --replicas=0


kubectl exec -it noobaa-db-pg-0 -- bash
>  psql -h 127.0.0.1 -p 5432 -U postgres
>  SELECT pg_terminate_backend (pid) FROM pg_stat_activity WHERE datname = 'nbcore';
  
kubectl cp mcg.bck noobaa-db-pg-0:/var/lib/pgsql/test.db
kubectl exec -it noobaa-db-pg-0 -- pg_restore -d nbcore /var/lib/pgsql/test.db -c
```


```bash
kubectl delete secret noobaa-root-master-key-backend; 
kubectl delete secret noobaa-root-master-key-volume;
kubectl delete secret noobaa-admin; 
kubectl delete secret noobaa-db; 
kubectl delete secret noobaa-operator; 
kubectl delete secret noobaa-server; 
kubectl delete secret  noobaa-endpoints; 


kubectl apply -f ./yamls/noobaa-root-master-key-backend.yaml
kubectl apply -f ./yamls/noobaa-root-master-key-volume.yaml
kubectl apply -f ./yamls/noobaa-admin.yaml
kubectl apply -f ./yamls/noobaa-db.yaml
kubectl apply -f ./yamls/noobaa-operator.yaml
kubectl apply -f ./yamls/noobaa-server.yaml
kubectl apply -f ./yamls/noobaa-endpoints.yaml
```


```bash
kubectl scale deployment noobaa-operator --replicas=1
kubectl scale deployment noobaa-endpoint --replicas=1
kubectl scale sts noobaa-core --replicas=1
```

```bash
kubectl port-forward -n noobaa service/s3 10443:443

export NOOBAA_ACCESS_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d')
export NOOBAA_SECRET_KEY=$(kubectl get secret noobaa-admin -n noobaa -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d')
alias s3='AWS_ACCESS_KEY_ID=$NOOBAA_ACCESS_KEY AWS_SECRET_ACCESS_KEY=$NOOBAA_SECRET_KEY aws --endpoint https://localhost:10443 --no-verify-ssl s3'
export AWS_REQUEST_CHECKSUM_CALCULATION=when_required
export AWS_RESPONSE_CHECKSUM_CALCULATION=when_required

s3 ls 
```