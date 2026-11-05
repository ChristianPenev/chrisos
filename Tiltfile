# PostgreSQL
k8s_yaml('./infrastructure/k8s/postgres/postgres.yaml')

# User Service
docker_build(
    'chrisos/user-service',
    './services/user',
)

k8s_yaml([
    './infrastructure/k8s/user/deployment.yaml',
    './infrastructure/k8s/user/service.yaml',
])

k8s_resource(
    'user-service',
    port_forwards=8080,
)