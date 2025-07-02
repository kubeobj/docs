
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
kubectl apply -f ./yamls/aws-s3-secret.yaml
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

/* ERROR WHILE RESTORE TO DIFFERENT CLUSTER */

```bash
Jun-24 8:41:08.667 [Upgrade/25]    [L0] core.server.system_services.system_store:: system_store is running in standalone mode. skip _register_for_changes
Jun-24 8:41:08.668 [Upgrade/25] [ERROR] UPGRADE:: failed to load system store!! Error: NO_SUCH_KEY
    at MasterKeysManager.get_master_key_by_id (/root/node_modules/noobaa-core/src/server/system_services/master_key_manager.js:165:35)
    at MasterKeysManager._resolve_master_key (/root/node_modules/noobaa-core/src/server/system_services/master_key_manager.js:185:32)
    at MasterKeysManager.get_master_key_by_id (/root/node_modules/noobaa-core/src/server/system_services/master_key_manager.js:167:27)
    at MasterKeysManager._resolve_master_key (/root/node_modules/noobaa-core/src/server/system_services/master_key_manager.js:185:32)
    at MasterKeysManager.get_master_key_by_id (/root/node_modules/noobaa-core/src/server/system_services/master_key_manager.js:167:27)
    at MasterKeysManager.decrypt_all_accounts_secret_keys (/root/node_modules/noobaa-core/src/server/system_services/master_key_manager.js:312:36)
    at /root/node_modules/noobaa-core/src/server/system_services/system_store.js:441:51
    at process.processTicksAndRejections (node:internal/process/task_queues:95:5)
    at async Semaphore.surround (/root/node_modules/noobaa-core/src/util/semaphore.js:71:84)
    at async init_db_upgrade (/root/node_modules/noobaa-core/src/upgrade/upgrade_manager.js:42:9)
```


**Noobaa Core CodeBase**
```javascript
    _resolve_master_key(m_key) {
        const { NOOBAA_ROOT_SECRET } = process.env;
        // in case we are resolving an old structured root-key, we will update to the new format
        if (!NOOBAA_ROOT_SECRET && this.is_root_key(m_key.master_key_id) && !m_key.root_key_id) {
            m_key.root_key_id = this.get_root_key_id();
            m_key.master_key_id = undefined;
        }
        // m_key.master_key_id._id doesn't exist when encrypting account secret keys and the account 
        // not yet inserted to db (in create account) or when master_key_id is the ROOT_KEY
        const m_of_mkey_id = (m_key.master_key_id && m_key.master_key_id._id) || m_key.master_key_id;
        const m_of_mkey = this.get_master_key_by_id(m_key.root_key_id || m_of_mkey_id || ROOT_KEY);
        if (!m_of_mkey) throw new Error('NO_SUCH_KEY');
    
        const iv = m_of_mkey.cipher_iv || Buffer.alloc(16);
        const decipher = crypto.createDecipheriv(m_of_mkey.cipher_type, m_of_mkey.cipher_key, iv);
        // cipher_key is Buffer and after load system - cipher_key is binary.
        const data = Buffer.isBuffer(m_key.cipher_key) ? m_key.cipher_key : Buffer.from(m_key.cipher_key.buffer, 'base64');
        let cipher_key = decipher.update(data);
        if (m_key.cipher_type !== 'aes-256-gcm') cipher_key = Buffer.concat([cipher_key, decipher.final()]);
        const decrypted_master_key = _.defaults({ cipher_key }, m_key);
    
        if (!Buffer.isBuffer(decrypted_master_key.cipher_iv)) {
            // we would like to keep it as Buffer in resolved_master_keys_by_id
            decrypted_master_key.cipher_iv = Buffer.from(decrypted_master_key.cipher_iv.buffer, 'base64');
        }
        this.resolved_master_keys_by_id[m_key._id.toString()] = decrypted_master_key;
        return decrypted_master_key;
}
```
Notes: 

- Seems the `master_key_id` is different while restoring to different cluster.
- So, Haven't finding the `master_key_id` and getting error `NO_SUCH_KEY`. 

