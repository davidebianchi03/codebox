FROM debian:bookworm

ARG USER_NAME
ARG USER_EMAIL=""
ARG USER_FIRST_NAME=""
ARG USER_LAST_NAME=""

# Install dependecies
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
    wget \
    build-essential \
    zip \
    unzip \
    python3 \
    python3-pip \
    python3-dev \
    python3-venv \
    && apt-get clean &&  \
    rm -rf /var/lib/apt/lists/*

# Create the user
RUN useradd -rm -d /home/${USER_NAME} -s /bin/bash -G sudo ${USER_NAME} && \
    echo "${USER_NAME}:password" | chpasswd
    
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

USER ${USER_NAME}

ENV HOME=/home/${USER_NAME}

# Install golang
RUN cd /tmp && \
    wget https://go.dev/dl/go1.23.1.linux-amd64.tar.gz -O go.tar.gz && \
    sudo tar -xzvf go.tar.gz -C /usr/local && \
    echo export PATH=$HOME/go/bin:/usr/local/go/bin:$PATH >> ~/.profile && \
    rm -rf /tmp/go.tar.gz

ENV PATH=$PATH:$HOME/go/bin:/usr/local/go/bin

# Install dlv and set permissions
RUN go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest
