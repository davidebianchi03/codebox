# Server configuration
Here is a list of the configuration parameters to use to customize a selfhosted codebox instance.

**Required parameters**
- `CODEBOX_EXTERNAL_URL`: the url of the codebox instance
- `CODEBOX_WILDCARD_DOMAIN`: codebox allows you to expose ports running HTTP-based services either public or with password authentication. The ports will be exposed through subdomains of this domain. You will need to define a DNS record with a name such as `*.codebox.my-domain.com`.

**Optional parameters**
- `CODEBOX_DEBUG`: `default: false` run server in debug mode, more logs about error will be printed, do not turn on this setting in production environment.
- `CODEBOX_DB_DRIVER`: `default: mysql` db driver to usen
- `CODEBOX_DB_HOST`: `default: localhost` the hostname/ip address where the dbms is running
- `CODEBOX_DB_PORT`: `default: 3306` the port where dbms is listening on
- `CODEBOX_DB_NAME`: `default: codebox` the name of the database
- `CODEBOX_DB_USER`: `default: codebox` the user used for authentication on dbms
- `CODEBOX_DB_PASSWORD`: `default: password` the password used for authentication on dbms
- `CODEBOX_SERVER_PORT`: `default: 8100` the port where codebox server is listening
- `CODEBOX_BG_TASKS_CONCURRENCY`: `default: 5` the concurrency of the background tasks
- `CODEBOX_REDIS_HOST`: `default: localhost` specifies the hostname or IP address of the Redis server.
- `CODEBOX_REDIS_PORT` `default: 6379` specifies the port number on which the Redis server is listening for connections
- `CODEBOX_DATA_PATH`: `default: ./data` the path where data is stored
- `CODEBOX_AUTH_COOKIE_NAME`: `default: codebox_auth_token` the name of the authentication cookie
- `CODEBOX_SUBDOMAIN_AUTH_COOKIE_NAME` `default: subdomain_codebox_auth_token` the name of the authentication cookie use in subdomains
