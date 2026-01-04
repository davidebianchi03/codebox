# CLI

The Codebox CLI is a client-side component that must be installed on users’ machines. It provides connectivity to remote workspaces via SSH and includes additional commands for managing workspaces.

## Installation

The CLI is available for most operating systems and CPU architectures. You can download it directly from your Codebox instance by navigating to the **CLI** section in the user menu located in the top-right corner of the interface.

Installation methods vary depending on the operating system and package manager used.

### Homebrew (macOS and Linux)

Homebrew is the recommended installation method for both **macOS** and **Linux** systems where Homebrew is available.

If Homebrew is not already installed, follow the official installation guide:

* [https://brew.sh/](https://brew.sh/)

#### Install via Homebrew

First, add the Codebox tap to your Homebrew sources:

```bash
brew tap codebox/codebox-cli https://gitlab.com/codebox4073715/codebox-homebrew-tap.git
brew update
```

If `codebox-cli` is already installed, uninstall it before proceeding:

```bash
brew uninstall codebox-cli
```

Then install the Codebox CLI:

```bash
brew install codebox/codebox-cli/codebox-cli
```

### Windows

On Windows, you can install the CLI by downloading the provided installer. The installer installs the Codebox CLI alongside other applications on your system.

Alternatively, you may download the standalone binaries if you prefer a manual installation.

### Linux (manual installation)

On Linux systems where Homebrew is not available or not desired, you can install the CLI using a distribution-specific package or standalone binaries.

#### Debian-based distributions

For Debian-based distributions, install the CLI using the provided `.deb` package:

```bash
dpkg -i codebox-cli-<arch>.deb
```

#### Other distributions

If your distribution is not supported, download the standalone binaries and place them in a directory included in your `PATH` (for example, `/usr/local/bin`).

### Manual installation (all platforms)

If you prefer not to use a package manager, you can install the CLI by downloading the standalone binaries from your Codebox instance.

## Quick Start

Before using the CLI, you must initialize the local environment. This command adds the required entry to your SSH configuration file.

```bash
codebox-cli init
```

To authenticate with your Codebox instance, run the `login` command. You will be prompted to enter the server URL and paste the generated access token.

```bash
codebox-cli login
```

```{note}
At the moment, the CLI supports being logged in to only one Codebox instance at a time.
```
