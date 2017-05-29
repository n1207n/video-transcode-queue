helm install -f task_queue/helm/queue_storage_config.yml --name=queue-storage stable/redis

helm install -f database/helm/database_config.yml --name=app-database stable/postgresql
