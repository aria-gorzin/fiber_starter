# Comprehensive VPS Setup Guide

This guide provides a comprehensive set of steps for setting up a production-ready VPS. It merges instructions for initial setup, security hardening, web server configuration, containerization, and monitoring.

## 1. Initial Server Setup

### 1.1. Create a New User with Sudo Permissions

First, log in as root and create a new non-root user with `sudo` privileges.

```sh
# Log in as root
ssh root@your-server-ip

# Create a new user
adduser newuser

# Add the user to the sudo group
usermod -aG sudo newuser

# Test the new user's sudo access
su - newuser
sudo apt update
```

### 1.2. Set Up SSH Key Authentication

Disable password-based login and use SSH keys for better security.

```sh
# On your local machine, generate an SSH key pair if you don’t have one
ssh-keygen -t ed25519 -C "your_email@example.com"

# Copy the SSH key to the new user on the server
ssh-copy-id -i ~/.ssh/id_ed25519.pub newuser@your-server-ip

# Test key-based login
ssh newuser@your-server-ip
```

## 2. Harden OpenSSH

Secure your SSH server by disabling root login and password authentication.

```sh
# Open SSH configuration file
sudo nano /etc/ssh/sshd_config
```

Modify the following lines in the file:

```
PermitRootLogin no
PasswordAuthentication no
UsePAM no
```

Restart the SSH service to apply the changes:

```sh
sudo systemctl restart ssh
```

**Important**: Test that you can still log in with your new user via SSH before logging out of your current session.

## 3. Set Up a Firewall (UFW)

Configure a firewall to control inbound and outbound traffic.

```sh
# Install UFW if not already installed
sudo apt install ufw

# Set default policies
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow necessary ports
sudo ufw allow OpenSSH    # Or your custom SSH port
sudo ufw allow 80/tcp     # HTTP
sudo ufw allow 443/tcp    # HTTPS
sudo ufw allow 'Nginx Full'

# Enable UFW
sudo ufw enable

# Check UFW status
sudo ufw status
```

## 4. (Optional) Install and Configure Fail2Ban

Protect your server from brute-force attacks by installing Fail2Ban.

```sh
# Install Fail2Ban
sudo apt install fail2ban

# Create a local configuration file
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local

# Edit the local configuration for SSH
sudo nano /etc/fail2ban/jail.local
```

Ensure the `[sshd]` section is enabled and configured:

```
[sshd]
enabled = true
port = 22 # Change this if you've modified your SSH port
maxretry = 5
bantime = 3600
```

Restart the Fail2Ban service:

```sh
sudo systemctl restart fail2ban

# Check Fail2Ban status
sudo fail2ban-client status
sudo fail2ban-client status sshd
```

## 5. Install Docker and Docker Compose

Set up Docker to run your application in containers.

1.  **Remove old versions**:
    ```sh
    for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done
    ```
2.  **Install dependencies**:
    ```sh
    sudo apt-get update
    sudo apt-get install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc
    ```
3.  **Add the repository to Apt sources**:
    ```sh
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update
    ```
4.  **Install Docker Engine**:
    ```sh
    sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    ```
5.  **Add your user to the `docker` group**:
    ```sh
    sudo groupadd docker
    sudo usermod -aG docker $USER
    newgrp docker
    ```
6.  **Configure Docker to start on boot**:
    ```sh
    sudo systemctl enable docker.service
    sudo systemctl enable containerd.service
    ```

## 6. Set Up Nginx and HTTPS with Certbot

Use Nginx as a reverse proxy and secure it with free TLS certificates from Let's Encrypt.

1.  **Install Nginx and Certbot**:
    ```sh
    sudo apt update
    sudo apt install nginx certbot python3-certbot-nginx
    ```
2.  **Create an Nginx configuration for your application** in `/etc/nginx/sites-available/your_domain.conf`:
    ```nginx
    server {
        listen 80;
        server_name your_domain.com www.your_domain.com;

        location / {
            proxy_pass http://localhost:YOUR_APP_PORT; # Point to your app's port
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
    ```
3.  **Enable the configuration**:
    ```sh
    sudo ln -s /etc/nginx/sites-available/your_domain.conf /etc/nginx/sites-enabled/
    sudo nginx -t
    sudo systemctl restart nginx
    ```
4.  **Obtain an HTTPS Certificate**:
    ```sh
    sudo certbot --nginx -d your_domain.com -d www.your_domain.com
    ```
    Certbot will automatically update your Nginx configuration for HTTPS.

5.  **Test Auto-Renewal**: Certbot sets up a cron job for automatic renewal. You can test it with a dry run:
    ```sh
    sudo certbot renew --dry-run
    ```

## 7. Automated Deployment with GitHub Actions

Automate your deployment process using a GitHub Actions workflow.

1.  **Create GitHub Secrets** in your repository settings:
    *   `DOCKER_USERNAME`: Your Docker Hub username.
    *   `DOCKER_PASSWORD`: Your Docker Hub password or access token.
    *   `VPS_HOST`: The IP address or hostname of your VPS.
    *   `VPS_USERNAME`: The SSH user for your VPS (e.g., `newuser`).
    *   `VPS_SSH_KEY`: Your private SSH key for accessing the VPS.

2.  **Create the workflow file** at `.github/workflows/deploy.yml`. This workflow will build and push a Docker image on changes to `main` and then deploy it to your VPS.

### Deploying your application manually to the cloud

First, build your image, e.g.: `docker build -t myapp .`.
If your cloud uses a different CPU architecture than your development
machine (e.g., you are on a Mac M1 and your cloud provider is amd64),
you'll want to build the image for that platform, e.g.:
`docker build --platform=linux/amd64 -t myapp .`.

Then, push it to your registry, e.g. `docker push myregistry.com/myapp`.

## 8. Secure Networking with Tailscale

For enhanced security, you can restrict SSH access to your Tailscale private network.

1.  **Install and configure Tailscale on the server**:
    ```sh
    sudo tailscale up --advertise-connector --ssh --advertise-tags=tag:db-egress,tag:server
    ```
2.  **Route traffic through the server** (run this on your development machine to access services like a database privately):
    ```sh
    sudo tailscale set --accept-routes=true
    ```

## 9. Monitoring with Prometheus, Grafana & Elasticsearch

Set up a monitoring stack using Docker Compose.

1.  **Prometheus**: For metrics collection.
    ```yaml
    # docker-compose.yml
    prometheus:
      image: prom/prometheus
      volumes:
        - ./prometheus.yml:/etc/prometheus/prometheus.yml
      ports:
        - "9090:9090"
    ```
2.  **Grafana**: For visualizing metrics.
    ```yaml
    # docker-compose.yml
    grafana:
      image: grafana/grafana
      ports:
        - "3000:3000"
    ```
3.  **Elasticsearch & Kibana**: For logging.
    ```yaml
    # docker-compose.yml
    elasticsearch:
      image: docker.elastic.co/elasticsearch/elasticsearch:7.10.0
      environment:
        - discovery.type=single-node
      ports:
        - "9200:9200"

    kibana:
      image: docker.elastic.co/kibana/kibana:7.10.0
      ports:
        - "5601:5601"
    ```
4.  **Filebeat**: Install on the host to ship logs to Elasticsearch.
    ```sh
    curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-7.10.0-amd64.deb
    sudo dpkg -i filebeat-7.10.0-amd64.deb
    ```
    Configure `/etc/filebeat/filebeat.yml` to point to your Elasticsearch instance.

## 10. Advanced Configurations

### 10.1. Load Balancing with Nginx

If you have multiple application instances, you can use Nginx as a load balancer.

```nginx
upstream backend {
    server app1.your_domain.com;
    server app2.your_domain.com;
}

server {
    listen 80;
    server_name your_domain.com;

    location / {
        proxy_pass http://backend;
    }
}
```

### 10.2. Securely Manage Passwords with Docker Secrets

Avoid hardcoding passwords. Use Docker Secrets for sensitive data like database credentials.

1.  **Create a Docker Secret**:
    ```sh
    echo "your_postgres_password" | docker secret create postgres_password -
    ```
2.  **Use the secret in your `docker-compose.yml`**:
    ```yaml
    services:
      postgres:
        image: postgres
        environment:
          POSTGRES_PASSWORD_FILE: /run/secrets/postgres_password
        secrets:
          - postgres_password

    secrets:
      postgres_password:
        external: true
    ```
