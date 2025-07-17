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
