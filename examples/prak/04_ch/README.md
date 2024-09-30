```bash
CLUSTER_HOST=`terraform output cluster_endpoint`
DB_NAME=`terraform output db_name`
DB_PASSWORD=`terraform output db_password`
TOKEN=$(yc iam create-token)
docker run --rm cr.yandex/sol/edu-checker validate clickhouse --token $TOKEN\
  --host $CLUSTER_HOST\
  --database $DB_NAME\
  --username user\
  --password $DB_PASSWORD
```
