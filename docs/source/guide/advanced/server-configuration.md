# Server configuration

Here is a list of the configuration parameters to use to customize a selfhosted codebox instance.

## Required parameters

This params are required to deploy codebox server.

### CODEBOX_EXTERNAL_URL

This is the public url of the codebox instance. Is the url where users can view the web app, connect with the cli.

Example:

```bash
CODEBOX_EXTERNAL_URL=https://codebox.my-domain.com
```

### CODEBOX_WILDCARD_DOMAIN

Codebox allows you to expose ports running HTTP-based services either public or with password authentication. The ports will be exposed through subdomains of this domain. You will need to define a DNS record with a name such as `*.codebox.my-domain.com`.

Example:

```bash
CODEBOX_WILDCARD_DOMAIN=codebox.my-domain.com
```

## Email Sender

The email sender is optional but recommended. Certain features, such as user sign-up, are disabled if the email sender is not configured. Additionally, the email sender is used to deliver security notifications to administrators. The email sender must be an SMTP server.

### CODEBOX_EMAIL_SMTP_HOST

The hostname of the SMTP server.

Example:

```bash
CODEBOX_EMAIL_SMTP_HOST=mail.my-domain.com
```

### CODEBOX_EMAIL_SMTP_PORT

The port on which the SMTP server is available. Both SSL and non-SSL connections are supported.

Example:

```bash
CODEBOX_EMAIL_SMTP_PORT=465
```

### CODEBOX_EMAIL_SMTP_USER

The username used to authenticate to the SMTP server. This is also the sender of the emails.

Example:

```bash
CODEBOX_EMAIL_SMTP_USER=codebox@my-domain.com
```

### CODEBOX_EMAIL_SMTP_PASSWORD

The password used to authenticate with the SMTP server.

Example:

```bash
CODEBOX_EMAIL_SMTP_PASSWORD=password
```

## Advanced configuration

```{warning}
Changing these parameters requires caution, as incorrect values can break the Codebox installation.
```

### CODEBOX_BG_TASKS_CONCURRENCY

This is the concurrency of the background tasks, consider to increase this value in large installations. The default is 5.

```bash
CODEBOX_BG_TASKS_CONCURRENCY=5
```

### CODEBOX_USE_SUBDOMAINS

Codebox allows to expose services using subdomains. If you don't want to use them, turn off this setting. The services will be exposed with sub-urls, you may have to configure the exposed services to accept the codebox url as prefix.

```bash
CODEBOX_USE_SUBDOMAINS=false
```

### CODEBOX_AUTH_COOKIE_NAME

This is the name of the cookie where is stored authentication token on codebox

```bash
CODEBOX_AUTH_COOKIE_NAME=codebox_auth_token
```

### CODEBOX_SUBDOMAIN_AUTH_COOKIE_NAME

This is the name of the cookie where is stored authentication token on subdomains. The name of the cookie is the same for all subdomains, consider to change it if it clashes with a service cookie.

```bash
CODEBOX_SUBDOMAIN_AUTH_COOKIE_NAME=subdomain_codebox_auth_token
```
