# Server configuration
Here is a list of the configuration parameters to use to customize a selfhosted codebox instance.

- `CODEBOX_EXTERNAL_URL`: the url of the codebox instance
- `CODEBOX_WILDCARD_DOMAIN`: codebox allows you to expose ports running HTTP-based services either public or with password authentication. The ports will be exposed through subdomains of this domain. You will need to define a DNS record with a name such as `*.codebox.my-domain.com`.
