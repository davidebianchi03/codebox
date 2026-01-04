# Server CLI

The Codebox server provides several CLI commands for server management and disaster recovery. Below is a list of available commands:

## runserver

Starts the Codebox server. This is the default command for Codebox Docker images.

```bash
codebox runserver
```

## set-password

Sets or resets the password for a user. Use this command to reset a password if a user has lost it.

```bash
codebox set-password
```

## reset-ratelimit

Clears all rate limits applied to endpoints.

```bash
codebox reset-ratelimit
```
