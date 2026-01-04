# Devcontainer
Codebox supports the DevContainers standard for defining the structure of a workspace.

## Devcontainer files
Codebox uses the DevContainers specification to configure workspaces. The configuration files must be loaded from a Git repository, template-based loading is not supported for DevContainers.

When initializing the workspace, Codebox parses the devcontainer.json file and looks for the following keys:

`workspaceFolder` – Specifies the directory path where the Git repository will be automatically cloned. A persistent volume will be mounted at this location.
`remoteUser` – Indicates the default user to be used inside the container.

⚠️ Note: The `workspaceMount` key is not supported and should not be used.

Codebox supports DevContainer configurations using either a single container or a multi-container setup defined through a Docker Compose file.

In the case you are using a multi-container setup, you can use the same labels available for [Docker Compose based workspaces](./docker-compose.md) to customize the stack.

⚠️ Note: Git is required in the main container.

## Environment
Codebox provides a set of default environment variables that can be used within your Docker Compose configuration. All of these variables — except for the email address — are automatically converted to lowercase before being injected. You can also define your own custom environment variables as needed. The default environment variables include:
- `CODEBOX_WORKSPACE_ID` – The unique identifier of the workspace
- `CODEBOX_WORKSPACE_NAME` – The name of the workspace
- `CODEBOX_WORKSPACE_OWNER_EMAIL` – The email address of the workspace owner
- `CODEBOX_WORKSPACE_OWNER_FIRST_NAME` – The first name of the workspace owner
- `CODEBOX_WORKSPACE_OWNER_LAST_NAME` – The last name of the workspace owner
- `CODEBOX_WORKSPACE_RUNNER_ID` – The ID of the runner managing the workspace
- `CODEBOX_WORKSPACE_RUNNER_NAME` – The name of the runner managing the workspace

### Example
```json
{
	"name": "001-single-container",
	"image": "mcr.microsoft.com/devcontainers/base:latest",
	"workspaceFolder": "/workspace",
	"customizations": {
		"vscode": {
			"extensions": []
		},
		"settings": {}
	},
	"remoteUser": "vscode"
}
```
You can view the full source code and more examples [here](https://gitlab.com/codebox4073715/codebox/-/tree/master/examples/devcontainer). 
