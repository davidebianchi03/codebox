FROM nginx:1.27.2

ENV USERNAME=codebox

# Install docker and redis
RUN apt-get update && apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    iptables \
    --no-install-recommends && \
    curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list && \
    apt-get update && \
    apt-get install -y \
    nano \
    sudo \
    docker-ce \
    docker-ce-cli \
    containerd.io \
    docker-buildx-plugin \
    docker-compose-plugin \
    && apt-get clean &&  \
    rm -rf /var/lib/apt/lists/*

COPY --chown=${USERNAME}:${USERNAME} ./bin/ /codebox/bin
COPY --chown=${USERNAME}:${USERNAME} ./docker/production.env /codebox/bin/codebox.env
COPY --chown=${USERNAME}:${USERNAME} ./app/build /codebox/frontend
COPY --chown=${USERNAME}:${USERNAME} ./docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY --chown=${USERNAME}:${USERNAME} ./docker/startgin.sh /docker-entrypoint.d/

RUN chmod +x /docker-entrypoint.d/startgin.sh && \
    mkdir -p /codebox/bin/migrations

EXPOSE 8000

VOLUME [ "/codebox/db" ]
VOLUME [ "/codebox/data" ]
