# PostgreSQL
k8s_yaml('./infrastructure/k8s/postgres/postgres.yaml')

# User Service
docker_build(
    'chrisos/user-service',
    '.',  # repo root, so the shared ./pkg and ./gen/go modules are in the build context
    dockerfile='./services/user-service/Dockerfile',
    only=['./pkg', './gen/go', './services/user-service'],
)

k8s_yaml([
    './infrastructure/k8s/user-service/deployment.yaml',
    './infrastructure/k8s/user-service/service.yaml',
])

k8s_resource(
    'user-service',
    port_forwards=9090,
    resource_deps=['postgres'],
)

# People API
docker_build(
    'chrisos/people-api',
    '.',
    dockerfile='./apps/people/api/Dockerfile',
    only=['./pkg', './gen/go', './apps/people/api'],
)

k8s_yaml([
    './infrastructure/k8s/people-api/deployment.yaml',
    './infrastructure/k8s/people-api/service.yaml',
])

k8s_resource(
    'people-api',
    port_forwards=8080,
    resource_deps=['user-service'],
)
