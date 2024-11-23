FROM nginx:1.27.2

# Install docker
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

# Install NVM
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/master/install.sh | bash

# Load NVM into the current shell session
RUN export NVM_DIR="$HOME/.nvm" && \
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && \
    [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion" && \
    nvm install 20.12.2 && \
    nvm alias default 20.12.2

# Install devcontainers CLI
RUN npm install -g @devcontainers/cli

COPY ./bin/ /codebox/bin
COPY ./docker/production.env /codebox/bin/codebox.env
COPY ./app/build /codebox/frontend
COPY ./docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY ./docker/startgin.sh /docker-entrypoint.d/

RUN chmod +x /docker-entrypoint.d/startgin.sh && \
    chmod +x /codebox/bin/codebox && \
    chmod +x /codebox/bin/codebox-cli-linux-amd64 && \
    chmod +x /codebox/bin/agent && \
    mkdir -p /codebox/bin/migrations

EXPOSE 8000

VOLUME [ "/codebox/db" ]
VOLUME [ "/codebox/data" ]
