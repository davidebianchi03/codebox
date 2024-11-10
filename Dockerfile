FROM nginx:1.27.2

ENV USERNAME=codebox

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

# Create new user
RUN adduser --disabled-password --gecos '' ${USERNAME} && \
    usermod -a -G sudo ${USERNAME} && \
    usermod -a -G docker ${USERNAME} && \
    echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> \
    /etc/sudoers

USER ${USERNAME}

COPY --chown=${USERNAME}:${USERNAME} ./bin/ /codebox/bin
COPY --chown=${USERNAME}:${USERNAME} ./app/build /codebox/frontend
