FROM nginx:1.27.2

# Install docker
RUN apt-get update && apt-get install -y \
    apt-transport-https \
    ca-certificates \
    wget \
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

ENV HOME=/root

# Install golang
RUN cd /tmp && \
    wget https://go.dev/dl/go1.23.1.linux-amd64.tar.gz -O go.tar.gz && \
    tar -xzvf go.tar.gz -C /usr/local && \
    echo export PATH=$HOME/go/bin:/usr/local/go/bin:$PATH >> ~/.profile && \
    rm -rf /tmp/go.tar.gz

ENV PATH=$PATH:$HOME/go/bin:/usr/local/go/bin

# Install atlas
RUN curl -sSf https://atlasgo.sh | sh

COPY ./bin/ /codebox/bin
COPY ./docker/production.env /codebox/bin/codebox.env
COPY ./app/build /codebox/frontend
COPY ./docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY ./docker/startgin.sh /docker-entrypoint.d/
COPY ./atlas.hcl /codebox/bin
COPY ./migrations/ /codebox/bin/migrations
COPY ./templates/ /codebox/bin/templates

RUN chmod +x /docker-entrypoint.d/startgin.sh && \
    chmod +x /codebox/bin/codebox && \
    chmod +x /codebox/bin/codebox-cli-linux-amd64

EXPOSE 8000

WORKDIR /codebox

VOLUME [ "/codebox/data" ]
