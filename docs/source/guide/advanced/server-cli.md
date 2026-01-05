# Server CLI

The Codebox server exposes a set of CLI commands that can be used for server administration, user management, and disaster recovery. These commands are typically executed inside the Codebox server container or on the host where the Codebox binary is available.

Below is a detailed list of the available commands and their purpose.

## runserver

Starts the Codebox server process. This command launches all required services and is the **default command** used by the official Codebox Docker images. In most setups, this command does not need to be executed manually unless you are running Codebox outside of Docker or debugging startup behavior.

```bash
codebox runserver
```

## set-password

Sets or resets the password for an existing user. This command is useful for account recovery scenarios, such as when a user has lost access to their credentials or cannot log in.

The command will prompt for the required information (such as the target user and the new password).

```bash
codebox set-password
```

## reset-ratelimit

Clears all active rate limits and bans applied to IP addresses across all endpoints. This command should be used with caution and is typically intended for administrative intervention after false positives or during incident recovery.

```bash
codebox reset-ratelimit
```

## approve-user

Marks a user account as approved. This command is required when the **"Users must be approved"** setting is enabled and a user has completed sign-up and email verification but is still awaiting administrator approval.

The command requires the email address of the user to approve.

```bash
codebox approve-user --user-email user@mydomain.com
```

## verify-email

Manually marks a user’s email address as verified. This is useful if a user is unable to complete the email verification flow or if email delivery is temporarily unavailable.

The command requires the email address to be verified.

```bash
codebox verify-email --email-address user@mydomain.com
```
