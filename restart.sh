# ------------------------------------------ stop services
docker compose down

docker rmi luinc-pong-backend --force
docker rmi luinc-pong-frontend --force

# ------------------------------------------ create networks
docker network create traefik-public
docker network create internal-network

# ------------------------------------------ start services
# docker compose -f docker-compose.traefik.yml up -d
# docker compose -f docker-compose.yml up -d

docker compose up