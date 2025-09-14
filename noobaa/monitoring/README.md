```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install prometheus-stack prometheus-community/kube-prometheus-stack -n monitoring
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

```bash
kubectl label servicemonitor -n noobaa s3-service-monitor release=prometheus-stack
kubectl label servicemonitor -n noobaa noobaa-mgmt-service-monitor release=prometheus-stack
```
## To list down all the metrics:

```bash
kubectl get svc -n noobaa -l noobaa-mgmt-svc=true
kubectl get svc -n noobaa -l noobaa-s3-svc=true



kubectl port-forward -n noobaa svc/noobaa-mgmt 8080:80

curl http://localhost:8080/metrics/web_server
curl http://localhost:8080/metrics/bg_workers
curl http://localhost:8080/metrics/hosted_agents



kubectl port-forward -n noobaa svc/s3 8081:80
curl http://localhost:8081/



```