helm install -f task_queue/helm/queue-storage-config.yml --name=queue-storage stable/redis

helm install -f database/helm/database-config.yml --name=app-database stable/postgresql
