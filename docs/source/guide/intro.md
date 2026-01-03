# Intro

## What is codebox?

Codebox is a self-hosted distributed provider of remote development environments.

What does that mean? With Codebox you can define the resources and the structure of a workspace using standard formats like Docker Compose or Devcontainers. Moreover it provides connection to the workspaces through an SSH connection and the possibility to expose to everyone or with restrictions HTTP services running inside the containers.

Codebox consists in four main parts:

### 1. A central server with web UI

This is where you can view, create and manage workspaces. The UI provides also an editor for workspace templates and admin tools to manage the server and the connected services.

### 2. Runners that host and manage the workspaces

Here’s where Codebox’s architecture becomes flexible. Workspaces are not managed directly from the server. Instead, you register runners, each capable of running different type of running one or more workspace types. This leads to several benefits including the fact that the workload can be divided between multiple machines.

Runners must be able to reach the main server, but not the vice versa. In this way you can register runners from multiple locations without opening ports on routers.

#### 3. Agents running inside containers

Agents are running inside workspaces, they provide the connections to the containers. They have an integrated SSH server. The SSH connection is tunneled over Web Sockets, in this way you don’t need to open other ports on your router.

### 4. A CLI for connecting via SSH

The CLI is a component to install on users’ PCs. It provides an SSH proxy to connect via an SSH connection to the workspaces. You can also use the official VS Code extension that wraps the CLI and provides an easy way to connect to workspaces.

## Sources

The source code it's available here:

#### Main Server

- [https://github.com/davidebianchi03/codebox](https://github.com/davidebianchi03/codebox)
- [https://gitlab.com/codebox4073715/codebox](https://gitlab.com/codebox4073715/codebox)

#### Runner

- [https://gitlab.com/codebox4073715/codebox-docker-runner](https://gitlab.com/codebox4073715/codebox-docker-runner)

#### Agent

- [https://gitlab.com/codebox4073715/codebox-agent](https://gitlab.com/codebox4073715/codebox-agent)


#### CLI

- [https://gitlab.com/codebox4073715/codebox-cli](https://gitlab.com/codebox4073715/codebox-cli)

#### VS Code Extension

- [https://gitlab.com/codebox4073715/codebox-vscode-extension](https://gitlab.com/codebox4073715/codebox-vscode-extension)
