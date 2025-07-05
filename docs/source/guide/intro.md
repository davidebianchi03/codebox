# Intro

## What is codebox?

Codebox is a tool for provisioning remote development environments. What does that mean? It allows you to define the structure of workspaces composed of containers or virtual machines (VMs). You can connect to these remote containers or VMs via a tunneled SSH connection provided by Codebox. Additionally, Codebox offers extensions for some IDEs to simplify the connection process.

## How does codebox works?

Codebox is composed of four main components:
- **Server**: Handles user authentication, provides a web UI to manage workspaces, and includes all administrative tools.
- **Runner**: Connects to the server and is responsible for creating, starting, and managing workspaces. When a workspace is started, the server sends a request to the Runner to launch it.
- **Agent**: A program that runs inside containers or virtual machines. It provides an SSH server and enables port forwarding. It will be automatically installed by the runner.
- **CLI**: A command-line tool that runs on the user's machine. It allows users to forward ports and connect to remote workspaces via SSH.
