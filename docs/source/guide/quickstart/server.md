# Server

This document describes the procedure for deploying a self-hosted instance of Codebox.

## Deploy the Server

The easiest way to deploy your Codebox instance is using the automated setup script. The script will handle the installation and configuration for you.

The installer will guide you through the setup process, including configuring the required settings. Download and run the installer script:

```bash
curl --output codebox-installer.sh "https://gitlab.com/api/v4/projects/68940432/packages/generic/codebox-installer/{version}/codebox-installer.sh"
chmod +x codebox-installer.sh
sudo ./codebox-installer.sh
```

To upgrade the service, use the same commands used for installation.

```{warning}
Docker and Docker Compose are required. Ensure they are installed and available in your system PATH before continuing.
```

To deploy Codebox, first download the `docker-compose.yml` file:

```bash
curl --output docker-compose.yml "[https://gitlab.com/api/v4/projects/68940432/repository/files/docker-compose.yml/raw?ref=master](https://gitlab.com/api/v4/projects/68940432/repository/files/docker-compose.yml/raw?ref=master)"
```

### Settings

The Docker stack requires some configuration before it can start correctly. If you're using the setup script, configure the required settings in the `.env` (e.g. `/etc/codebox/codebox.env`) file. Otherwise, when deploying the Docker stack directly, provide them as stack environment variables.

When using the setup script, the server URL, whether subdomains are enabled, and the subdomain configuration are automatically handled by the installation script.

#### Environment variables

* `CODEBOX_EXTERNAL_URL`: (e.g. `https://codebox.example.com`) The externally accessible URL of the Codebox instance (for example, `https://codebox.my-domain.com`).

* `CODEBOX_USE_SUBDOMAINS`: Enables the use of subdomains to expose services. If enabled, this requires `CODEBOX_WILDCARD_DOMAIN`.

* `CODEBOX_WILDCARD_DOMAIN`: (e.g. `codebox.example.com`) Required if Codebox exposes HTTP-based services through dynamically generated subdomains. It defines the wildcard domain used for this purpose. A DNS record such as `*.codebox.my-domain.com` must be configured and point to the Codebox host.


#### Optional: Email Sender Configuration

Configuring an email sender service is strongly recommended. The email sender must be an SMTP-compatible server. Certain features (such as user sign-up and account approval workflows) are disabled if email delivery is not configured.

The following environment variables are used to configure the SMTP service:

* `CODEBOX_EMAIL_SMTP_HOST`: Hostname of the SMTP server.
* `CODEBOX_EMAIL_SMTP_PORT`: Port on which the SMTP server is listening.
* `CODEBOX_EMAIL_SMTP_USER`: Username used to authenticate with the SMTP server.
* `CODEBOX_EMAIL_SMTP_PASSWORD`: Password used to authenticate with the SMTP server.

A complete list of supported configuration parameters is available in the [configuration](../advanced/server-configuration) documentation.

## Reverse Proxy Configuration

If the Codebox instance is exposed through an HTTP reverse proxy (for example, Nginx or Apache), WebSocket support **must** be enabled to ensure proper operation of interactive features.

Refer to the official documentation of your proxy server for configuration details:

* [Enable WebSockets with Nginx](https://nginx.org/en/docs/http/websocket.html)
* [Enable WebSockets with Apache](https://httpd.apache.org/docs/current/mod/mod_proxy_wstunnel.html)
