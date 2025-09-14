# ------------------------------------------ stop services
docker compose down

docker rmi app-frontend:latest --force
docker rmi app-backend:latest --force

# ------------------------------------------ create networks
docker network create traefik-public

# ------------------------------------------ start services
docker compose -f docker-compose.traefik.yml up -d
docker compose -f docker-compose.yml up -d

# docker compose up