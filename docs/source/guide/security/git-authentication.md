# GIT Authentication

Codebox generates an SSH key pair for each user. The public key is available under your profile details, you can add the public key to your Git server to enable authentication. This key pair is also used to authenticate requests for retrieving workspace configurations from Git repositories.

The private key is not injected into containers, codebox automatically exports a custom `GIT_SSH_COMMAND` in workspace containers.
