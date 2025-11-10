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
    build-essential \
    openvpn \
    zip \
    unzip \
    bridge-utils \
    libcap2-bin \
    python3 \
    python3-venv \
    python3-pip \
    python3-dev \
    net-tools \
    && apt-get clean &&  \
    rm -rf /var/lib/apt/lists/*

# Create tun device
RUN mkdir -p /dev/net && \
    mknod /dev/net/tun c 10 200 && \
    chmod 777 /dev/net/tun

# Create user "vscode"
RUN useradd -rm -d /home/vscode -s /bin/bash -G sudo vscode && \
    echo "vscode:password" | chpasswd
    
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN groupadd wireguard && \
    usermod -aG wireguard vscode && \
    echo '%wireguard ALL=NOPASSWD: SETENV: /usr/bin/wg-quick' >> /etc/sudoers

USER vscode

ENV HOME=/home/vscode

# Install golang
RUN cd /tmp && \
    wget https://go.dev/dl/go1.23.1.linux-amd64.tar.gz -O go.tar.gz && \
    sudo tar -xzvf go.tar.gz -C /usr/local && \
    echo export PATH=$HOME/go/bin:/usr/local/go/bin:$PATH >> ~/.profile && \
    rm -rf /tmp/go.tar.gz

ENV PATH=$PATH:$HOME/go/bin:/usr/local/go/bin

ENV NVM_DIR=$HOME/.nvm
ENV NODE_VERSION=20.12.2

# Install NVM
RUN mkdir -p ${NVM_DIR} && \
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/master/install.sh | bash && \
    \. $NVM_DIR/nvm.sh && \
    \. $NVM_DIR/bash_completion && \
    nvm install $NODE_VERSION && \
    nvm alias default $NODE_VERSION && \
    nvm use default

# install atlas
RUN mkdir -p /tmp/atlas && \
    cd /tmp/atlas && \
    wget https://atlasgo.sh -O atlas-installer.sh && \
    chmod +x atlas-installer.sh && \
    ./atlas-installer.sh -y && \
    cd ~ && \
    rm -rf /tmp/atlas

# add node and npm to path so the commands are available
ENV NODE_PATH=$NVM_DIR/v$NODE_VERSION/lib/node_modules
ENV PATH=$NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH
ENV PATH=$PATH:/home/node/.npm-global/bin

RUN npm install -g @devcontainers/cli

# Install dlv and set permissions
RUN go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    sudo setcap cap_net_admin=+ep ${HOME}/go/bin/dlv

# Prepare VSCode server directory and set permissions
# Create a symbolic link to persist VSCode server data
RUN mkdir -p $HOME/.vscode-server && \
    sudo mkdir -p /vscode-server && \
    sudo chown -R vscode:vscode $HOME/.vscode-server && \
    sudo chown -R vscode:vscode /vscode-server && \
    ln -s /vscode-server $HOME/.vscode-server

# Create a symbolic link to persist Go bin data
RUN mkdir -p $HOME/go && \
    sudo mkdir -p /go-home && \
    sudo chown -R vscode:vscode $HOME/go && \
    sudo chown -R vscode:vscode /go-home && \
    ln -s /go-home $HOME/go

ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

USER root


# sudo setcap cap_net_admin=eip /usr/local/go/bin/go
# export PATH=/home/vscode/go/bin:/usr/local/go/bin:$PATH
# sudo setcap cap_net_admin=+ep /home/vscode/go/bin/dlv
# sudo setcap cap_net_admin=+ep /home/vscode/.vscode-server/cli/servers/Stable-fabdb6a30b49f79a7aba0f2ad9df9b399473380f/server/bin/remote-cli/code

# export LC_ALL=C.UTF-8
# export LANG=C.UTF-8
