# ------------------------------------------ stop services
docker compose down

docker rmi app-frontend:latest --force
docker rmi ghcr.io/jda5/luinc-pong:latest --force
docker rmi luinc-pong-frontend --force
docker rmi luinc-pong-backend --force

# ------------------------------------------ create networks
docker network create traefik-public
docker network create internal-network

# ------------------------------------------ start services
docker compose -f docker-compose.traefik.yml up -d
docker compose -f docker-compose.yml up -d

# docker compose up