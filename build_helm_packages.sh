helm install -f task_queue/queue_storage_config.yml --name=queue-storage stable/redis

helm install -f database/database_config.yml --name=app-database stable/postgresql
