# PostgreSQL
k8s_yaml('./infrastructure/k8s/postgres/postgres.yaml')

# User Service
docker_build(
    'chrisos/user-service',
    '.',  # repo root, so the shared ./pkg module is in the build context
    dockerfile='./services/user-service/Dockerfile',
    only=['./pkg', './services/user-service'],
)

k8s_yaml([
    './infrastructure/k8s/user-service/deployment.yaml',
    './infrastructure/k8s/user-service/service.yaml',
])

k8s_resource(
    'user-service',
    port_forwards=8080,
)