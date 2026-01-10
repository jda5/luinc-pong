# LUinc. Pong 🏓

An Elo leaderboard for our office table tennis games.

## Development

### Prerequisites

- Docker and Docker Compose
- Node.js (for frontend development)
- Go (for backend development)

### Running Locally

1. Ensure the `build` section is uncommented in `docker-compose.yml`:
   ```yaml
   build:
     dockerfile: ./backend/Dockerfile
   ```

2. Comment out the `image` line:
   ```yaml
   # image: ghcr.io/jda5/luinc-pong:latest
   ```

3. Start the backend:
   ```bash
   docker compose up -d backend
   ```

4. The API will be available at `http://localhost:8080`

---

## Deployment

### 1. Authenticate with GitHub Container Registry

First, create a GitHub Personal Access Token:
1. Go to GitHub.com → Profile → **Settings**
2. **Developer settings** → **Personal access tokens** → **Tokens (classic)**
3. **Generate new token (classic)**
4. Select scopes: `write:packages`, `read:packages`
5. Copy the token

Then login to GHCR:
```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
echo $GITHUB_TOKEN | docker login ghcr.io -u jda5 --password-stdin
```

### 2. Build and Push the Docker Image

```bash
# Build using compose
docker compose build backend

# Tag the built image
docker tag luinc-pong-backend:latest ghcr.io/jda5/luinc-pong:latest

# Push to registry
docker push ghcr.io/jda5/luinc-pong:latest
```

### 3. Configure for Production

On your production server, modify `docker-compose.yml`:

1. **Comment out** the `build` section:
   ```yaml
   # build:
   #   dockerfile: ./backend/Dockerfile
   ```

2. **Uncomment** the `image` line:
   ```yaml
   image: ghcr.io/jda5/luinc-pong:latest
   ```

3. **Comment out** the `ports` section (Traefik handles routing):
   ```yaml
   # ports:
   #   - "8080:8080"
   ```

### 4. Deploy on Production Server

```bash
# Pull the latest image
docker compose pull backend

# Start/restart the service
docker compose up -d backend
```

---

## Environment Variables

Create a `.env` file at `./backend/src/.env` with the required database configuration.
