FROM ubuntu:22.04

ARG NODE_VERSION=22
ARG USERGROUP
ENV NVM_DIR=/usr/local/nvm
ENV PATH=$NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

RUN groupadd --force -g $USERGROUP docky && \
    useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP docky

RUN apt-get update && apt-get install -y \
    apt-utils \
    sudo \
    build-essential \
    curl \
    bash && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p $NVM_DIR && \
    curl -fsSL https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash && \
    export NVM_DIR=$NVM_DIR && \
    . "$NVM_DIR/nvm.sh" && \
    nvm install ${NODE_VERSION} && \
    nvm alias default ${NODE_VERSION} && \
    npm install -g pm2 && \
    ln -s "$NVM_DIR/versions/node/$(nvm version)/bin/node" /usr/bin/node && \
    ln -s "$NVM_DIR/versions/node/$(nvm version)/bin/npm" /usr/bin/npm && \
    ln -s "$NVM_DIR/versions/node/$(nvm version)/bin/npx" /usr/bin/npx && \
    ln -s "$NVM_DIR/versions/node/$(nvm version)/bin/pm2" /usr/bin/pm2 && \
    chmod -R a+rx "$NVM_DIR"

WORKDIR /var/www
