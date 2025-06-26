# Deployment Guide (Docker-Based)
This document explains how to deploy the application using Docker and docker-compose.

## Environment Configuration
The project uses a .env file for runtime configuration. This file is not committed to version control to avoid leaking credentials or sensitive information.

Instead, a .env.template file is provided with required environment variable keys. Before running the project, copy and rename this file:

```bash
cp docker/.env.template docker/.env
```
Then, edit docker/.env to include your environment-specific values.

### Example: docker/.env.template
```env

# pgadmin
PGADMIN_DEFAULT_EMAIL=
PGADMIN_DEFAULT_PASSWORD= 
PGADMIN_LISTEN_PORT=80

# PHARMACY data DB
DB_PHARMACY_DBNAME=postgres-pharmacy
DB_PHARMACY_USER=pharma_admin
DB_PHARMACY_PASSWORD=your_password

```
## PGAdmin Auto-Connect (Optional)
To auto-connect PGAdmin to your databases:

+ A servers.json file is mounted to the container.

+ This file is not committed to Git. Instead, a template is provided:

```bash
cp docker/servers.template.json docker/servers.json
```
Then fill in the appropriate connection details based on your environment.

## Step-by-Step Deployment

### 1. Build the Backend Image
```bash
docker compose -f docker-compose.yaml build
```
This command builds the phantom-be Go backend service into a Docker image using the Dockerfile.

### 2. Start the Databases
```bash
docker compose -f docker-compose-db.yaml up -d
``` 
This starts the required PostgreSQL containers:

- postgres-pharmacy: main application database

- postgres-user: optional user authentication DB

- pgadmin: for manual inspection

> Ensure these services are healthy before proceeding.

### 3. Run Schema Migrations & Load Sample Data
```bash
docker compose -f docker-compose-init.yaml up --exit-code-from init-preprocess
```
This performs:
- Schema migration (via migrateSchema)

- ETL preprocessing and data loading (initUsers, initPharmacies)

### 4. Start the Backend API Service
```bash
docker compose -f docker-compose.yaml up
```

This launches the Go API service (phantom-be) and exposes REST endpoints for external use

## Notes

+ Modify environment variables in docker/.env as needed.

+ Input data is located in the /data directory (mounted into the container).

+ Logs from init-preprocess indicate success/failure of initial loading.

To stop service containers:
```bash
docker compose -f docker-compose.yaml down
```

## Related Files
+ docker-compose.yaml: Backend service definition

+ docker-compose-db.yaml: PostgreSQL & pgAdmin configuration

+ docker-compose-init.yaml: Init containers for migrations and ETL