# Palm Recognition Backend - Server Deployment Guide

This guide provides step-by-step instructions on how to deploy the Go Backend and PostgreSQL Database to a production server (Ubuntu/Debian recommended).

---

## 1. Prerequisites

Before starting, ensure your server has the following installed:
* **Git:** To clone the repository.
* **Docker:** To run the containerized application.
* **Docker Compose:** To orchestrate the API and Database together.

```bash
# Ubuntu/Debian quick install commands
sudo apt update
sudo apt install git docker.io docker-compose -y
sudo systemctl enable --now docker
```

---

## 2. Clone the Repository

Clone your project onto the server and navigate into the backend directory:

```bash
git clone <your-repository-url>
cd palm-back-end
```

---

## 3. Configuration

You must create a `.env` file for your production environment. 

1. Copy the example file or create a new one:
   ```bash
   cp .env.example .env
   # OR
   nano .env
   ```

2. Configure your production variables. **CRITICAL:** Make sure you change the `JWT_SECRET` to a long, random string for security.

   ```env
   # Application Port
   APP_PORT=8080

   # PostgreSQL Database Configuration
   DB_HOST=db
   DB_PORT=5432
   DB_USER=postgres_production_user
   DB_PASSWORD=super_secure_db_password
   DB_NAME=palm_attendance

   # Security
   JWT_SECRET=YOUR_VERY_LONG_RANDOM_SECURE_STRING_HERE
   ```

---

## 4. Automatic Database Migrations

**You do not need to manually set up the database!**
Because of how our `docker-compose.yml` is structured, the `./migrations/001_init.sql` script is mounted directly into the Postgres container. 

The **very first time** the database container starts, it will automatically execute this script, enable UUID extensions, and create all your tables (users, devices, attendance_logs, etc).

---

## 5. Build and Deploy

To build the Go binary and start the containers in the background, run:

```bash
sudo docker-compose up -d --build
```

### Useful Docker Commands:
* **Check Status:** `sudo docker-compose ps` (Both `palm_api` and `palm_db` should say "Up").
* **View API Logs:** `sudo docker-compose logs -f api`
* **View DB Logs:** `sudo docker-compose logs -f db`
* **Stop Server:** `sudo docker-compose down`

---

## 6. Verification

Once deployed, you can verify the backend is running by sending a quick health check or fetching the API.

```bash
curl http://localhost:8080/api/v1/auth/login -X POST -H "Content-Type: application/json" -d '{}'
```
*(You should receive a JSON response indicating invalid request/credentials, which proves the API is alive and responding).*

---

## 7. Next Steps: Nginx & HTTPS (Recommended)

Since the physical Raspberry Pi devices and Mobile Apps will be sending sensitive data (passwords, encrypted palm templates) over the network, you **MUST** use HTTPS in production.

It is highly recommended to install **Nginx** and **Certbot (Let's Encrypt)** in front of this Docker container:

1. Map your domain name (e.g., `api.yourdomain.com`) to your server's IP address.
2. Install Nginx: `sudo apt install nginx`
3. Proxy traffic from port 80/443 down to your Docker container running on port 8080.
4. Run Certbot to generate a free SSL certificate.
