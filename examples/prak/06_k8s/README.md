```bash
CLUSTER_ID=`terraform output --raw k8s_cluster_id`
echo $CLUSTER_ID
```

```bash
yc managed-kubernetes cluster get-credentials $CLUSTER_ID  --external --force
```


Подготовить сертификаты для подключения к кластеру:
```bash
yc managed-kubernetes cluster get --id $CLUSTER_ID --format json | \
  jq -r .master.master_auth.cluster_ca_certificate | \
  awk '{gsub(/\\n/,"\n")}1' > ca.pem
```

Подготовит токен SA для подключения к кластеру:
```bash
SA_TOKEN=$(kubectl -n kube-system get secret $(kubectl -n kube-system get secret | \
  grep admin-user-token | \
  awk '{print $1}') -o json | \
  jq -r .data.token | \
  base64 -d)
```

```bash
MASTER_ENDPOINT=$(yc managed-kubernetes cluster get --id $CLUSTER_ID \
  --format json | \
  jq -r .master.endpoints.external_v4_endpoint)
echo $MASTER_ENDPOINT
```

```bash
kubectl config set-cluster yc-django-k8s-cluster \
  --certificate-authority=ca.pem \
  --server=$MASTER_ENDPOINT \
  --kubeconfig=test.kubeconfig
```

```bash
kubectl config set-credentials admin-user \
  --token=$SA_TOKEN \
  --kubeconfig=test.kubeconfig
```

```bash
kubectl config set-context default \
  --cluster=sa-test2 \
  --user=admin-user \
  --kubeconfig=test.kubeconfig
```

```bash
kubectl config use-context default \
  --kubeconfig=test.kubeconfig
```