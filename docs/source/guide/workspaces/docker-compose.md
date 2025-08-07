# Docker Compose
Codebox enables you to define the structure of a workspace using a standard docker-compose file. 

## Compose file
While you can use a regular Docker Compose configuration to define your workspace, Codebox provides a set of custom labels that can be added to containers to seamlessly integrate your stack with the Codebox environment.

The docker-compose file can be loaded from either a Git repository or a predefined template.

## Labels

### Expose a port
You can configure the default exposed ports for services using labels:
- Use the label `com.codebox.port.<service_name>` to bind a service to a specific port (e.g., `com.codebox.port.phpmyadmin=80`).
- Use the label `com.codebox.port.<service_name>.public` to define the port's visibility (e.g., `com.codebox.port.phpmyadmin.public=false`). By default, ports are not publicly accessible (`false`).

### Workspace path
You can specify the working directory path using the `com.codebox.workspace_path` label. This path will be used as the default location opened by IDE integrations and remote connections (e.g., `com.codebox.workspace_path=/workspace`).

### Container user
You can specify the default user for container access using the `com.codebox.user` label. If not set, Codebox will attempt to automatically determine the appropriate username (e.g., `com.codebox.user=user`).

## Environment variables
Codebox provides a set of default environment variables that can be used within your Docker Compose configuration. These variables are automatically converted to lowercase before being injected. You can also define your own custom environment variables. The default environment variables include:
- `CODEBOX_WORKSPACE_ID` – The unique identifier of the workspace
- `CODEBOX_WORKSPACE_NAME` – The name of the workspace
- `CODEBOX_WORKSPACE_OWNER_EMAIL` – The email address of the workspace owner
- `CODEBOX_WORKSPACE_OWNER_FIRST_NAME` – The first name of the workspace owner
- `CODEBOX_WORKSPACE_OWNER_LAST_NAME` – The last name of the workspace owner
- `CODEBOX_WORKSPACE_RUNNER_ID` – The ID of the runner managing the workspace
- `CODEBOX_WORKSPACE_RUNNER_NAME` – The name of the runner managing the workspace

### Example
```yml
version: "3.8"

services:
    dev:
        build:
            context: .
            args:
                USER_NAME: ${CODEBOX_WORKSPACE_OWNER_FIRST_NAME:-}
                USER_EMAIL: ${CODEBOX_WORKSPACE_OWNER_EMAIL:-}
                USER_FIRST_NAME: ${CODEBOX_WORKSPACE_OWNER_FIRST_NAME:-}
                USER_LAST_NAME: ${CODEBOX_WORKSPACE_OWNER_LAST_NAME:-}
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
You can view the full source code [here](https://gitlab.com/codebox4073715/codebox/examples/docker-compose/001-basic-stack).
