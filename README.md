# File Vault

Secure, deduplicated file storage platform (Go backend + React frontend).

## Stack
- **Backend**: Go, GraphQL, PostgreSQL, Redis, MinIO
- **Frontend**: Next.js + Tailwind (with shadcn/ui)
- **Dev**: Docker Compose


## Getting Started
### Prerequisites
- [Docker Compose](https://docs.docker.com/compose/) (for running containers)
- [Node.js & npm](https://nodejs.org/) (for managing frontend dependencies, dev builds, and scripts)

### Running Locally

1. **Clone the repository**
```bash
git clone https://github.com/iolynx/FileVault.git
cd FileVault
```
2. **Create a `.env` File in Project Root.** (Refer Configuration)

3. **Set up the Frontend**
```bash
cd frontend
npm install
npm run dev
```
The frontend will now be running at `http://localhost:3000`

4. **Run the backend**

Open a new terminal in the project root:
```bash
docker compose up --build
```
The backend services (API + PostgreSQL + MinIO + Redis) will start up inside Docker.

5. **Shut down services (when done)**
```bash
docker compose down
```

### Development
Use Air (watches backend and hot reloads go server on changes) for development 
```bash
docker compose up -d
cd backend && air
cd frontend && npm run dev
```

## Configuration
### Backend Configuration

Create a `.env` file in the project root and define the following variables.
These control database connections, storage, and other settings.

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `filevault` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `MINIO_ENDPOINT` | MinIO server address | `localhost:9000` |
| `MINIO_ACCESS` | MinIO access key | `minioadmin` |
| `MINIO_SECRET` | MinIO secret key | `minioadmin` |
| `MINIO_BUCKET` | MinIO bucket name | `filevault` |
| `MINIO_SECURE` | Use HTTPS for MinIO | `false` |
| `REDIS_ADDR` | Redis server address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password | `redis` |
| `REDIS_DB` | Redis database number | `0` |
| `PORT` | Backend server port | `8080` |
| `DEFAULT_STORAGE_QUOTA` | Default storage quota per user (bytes) | `10000000` |
| `API_RATE_LIMIT` | Max API requests per window | `2` |
| `API_RATE_LIMIT_WINDOW_SECONDS` | Rate limit window (seconds) | `1` |
| `JWT_SECRET` | Secret key for JWT tokens | `supersecret` |

> ⚠️ **Note:** After updating the `.env` file, make sure to restart the backend services so the changes take effect.

### Frontend Configuration

The frontend is a Next.js application. By default, it connects to the backend at `http://localhost:8080`.  

You can override this in `.env.local` at the frontend root if needed:

| Variable | Description | Example |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_BASE_URL` | Base URL for API requests | `http://localhost:8080` |

The frontend reads environment variables at build time, so any changes require restarting the development server.


Author: Vishal R
