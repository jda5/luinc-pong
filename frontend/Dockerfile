
# Use official Node.js LTS image
FROM node:20-alpine AS builder

WORKDIR /app/src

# Install dependencies
COPY frontend/src/package.json frontend/src/package-lock.json* ./
RUN npm install

# Copy source code
COPY frontend/src/ ./

# Build the app
RUN npm run build

# Production image
FROM node:20-alpine AS runner

WORKDIR /app/src

# Install only production dependencies
COPY frontend/src/package.json frontend/src/package-lock.json* ./
RUN npm install --omit=dev

# Copy built files from builder
COPY --from=builder /app/src/dist ./dist

COPY frontend/src/vite.config.ts .

EXPOSE 5173

CMD ["npx", "vite", "preview", "--port", "5173", "--host", "0.0.0.0"]