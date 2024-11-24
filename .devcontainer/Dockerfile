FROM nginx:1.27.2

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
    && apt-get clean &&  \
    rm -rf /var/lib/apt/lists/*

# Create user "vscode"
RUN useradd -rm -d /home/vscode -s /bin/bash -G sudo vscode && \
    echo "vscode:password" | chpasswd
    
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

USER vscode

# Install golang
RUN wget https://go.dev/dl/go1.23.1.linux-amd64.tar.gz -O go.tar.gz && \
    sudo tar -xzvf go.tar.gz -C /usr/local && \
    echo export PATH=$HOME/go/bin:/usr/local/go/bin:$PATH >> ~/.profile

ENV NVM_DIR=/home/vscode/.nvm
ENV NODE_VERSION=20.12.2

# Install NVM
RUN mkdir -p ${NVM_DIR} && \
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/master/install.sh | bash && \
    \. $NVM_DIR/nvm.sh && \
    \. $NVM_DIR/bash_completion && \
    nvm install $NODE_VERSION && \
    nvm alias default $NODE_VERSION && \
    nvm use default

# add node and npm to path so the commands are available
ENV NODE_PATH=$NVM_DIR/v$NODE_VERSION/lib/node_modules
ENV PATH=$NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

RUN npm install -g @devcontainers/cli

USER root
