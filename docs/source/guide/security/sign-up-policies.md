# Sign-Up Policies

By default, user sign-up is disabled. Only administrators can create new users.

Sign-up can be enabled from the `Authentication` section of the admin menu by selecting the **`Sign Up Open`** checkbox. To enable sign-up, the email sender must be configured.

Users who sign up receive a verification email after their first login attempt. They cannot log in until their email address has been verified. Administrators can manually verify user email addresses from the admin panel. Users created by administrators are automatically verified.

Some policies and restrictions can be configured to control sign-up behavior.

## Users Approval

Administrators can enable the **`Users must be approved`** setting. When enabled:

* New users cannot log in until approved by an administrator.
* Administrators can define regex patterns to automatically approve specific email domains or addresses.
* Administrators receive an email notification once a user has verified their email address.
* Users created by administrators are always verified.

If this setting was not initially enabled, administrators must manually verify all users who signed up and were not automatically verified.

**Examples of auto-approval regex:**

* `@example.com$` → automatically approve all emails ending with `@example.com`
* `^admin@.*` → automatically approve any email starting with `admin@`

## Sign-Up Restrictions

Administrators can restrict sign-up to only those email addresses that match one or more specified regex patterns by enabling the **`Sign Up Restricted`** setting. Allowed patterns must be specified in the **`Allowed Email Addresses Regex`** field.

**Examples of allowed email regex:**

* `@company.com$` → only allow emails ending with `@company.com`
* `^user[0-9]+@example.com$` → only allow emails like `user123@example.com`

## Sign-Up Blacklist

Administrators can define regex patterns to block specific email addresses. Any email that matches a blacklisted pattern is not allowed to sign up.

**Examples of blacklisted email regex:**

* `@spam.com$` → block all emails ending with `@spam.com`
* `^test@.*` → block any email starting with `test@`
