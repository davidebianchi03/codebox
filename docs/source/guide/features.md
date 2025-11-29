# Features

Here are listed the major features of Codebox.

## Distributed workspaces

Codebox workspaces are not directly managed by the central server; instead, you have to register runners. This decentralized approach allows you to distribute workspaces among multiple servers. Consequently, you can use different runners to access specific, segmented resources.

## Define workspace structure with standard formats

The structure of Codebox workspaces is defined using standard formats like Docker Compose or Devcontainers. While Codebox workspaces can be customized—for example, you can specify properties like the container username or workspace path using container labels in Docker Compose—the underlying standard remains.

e.g.
```yaml
version: "3.8"

services:
    dev:
        build:
            context: .
        stdin_open: true
        tty: true
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
            - workspace:/home/${CODEBOX_WORKSPACE_OWNER_FIRST_NAME:-}
        labels:
            - com.codebox.user=${CODEBOX_WORKSPACE_OWNER_FIRST_NAME:-}

volumes:
    workspace:
```

## Connect to workspace through an SSH connection

Codebox provides connection to the workspaces. You don't need to open SSH ports on your routers since the connection is tunneled over HTTP. To establish this connection, you will need to install a tool on your PC (the codebox-cli) that provides an SSH proxy. Codebox also provides a VS Code extension that wraps the CLI, making it even easier and more seamless to connect to and manage your workspaces.

## Expose HTTP ports

Codebox allows you to expose specific ports from the workspaces to the public internet, with or without access restrictions. Any services exposed this way will be accessible under subdomains of the Codebox server.

## In-browser terminal

Codebox provides an in-browser terminal for each workspace container. This can be useful for performing rapid operations.

Watch this space for new features...
