# Ratelimits

Some Codebox endpoints are protected by rate limits. Rate limiting is applied per IP address using an incremental strategy: if an IP exceeds the same limit multiple times, the ban interval increases based on the number of violations.

When a rate limit is exceeded multiple times, a notification is sent to all administrators, so they can take appropriate action if the activity is considered a potential violation.
