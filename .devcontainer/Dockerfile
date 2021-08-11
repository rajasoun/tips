# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.177.0/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version: 1, 1.16, 1.15
ARG VARIANT="1.16"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT} as base

COPY packages/*.* /tmp/packages-tmp/

# Install Node.js
FROM base as node-builder
ARG INSTALL_NODE="true"
ARG NODE_VERSION="lts/*"

RUN set -x \
    && CMD="umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1" \
    && if [ "${INSTALL_NODE}" = "true" ]; then su vscode -c "$CMD"; fi

# OS packages
FROM node-builder as os-builder
RUN set -x \
    && apt-get update \
    && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends bash bash-completion \
       vim git netcat jq curl git-flow wget g++ zsh-autosuggestions \
       zsh-syntax-highlighting binutils libcurl4 zlib1g libdw1 shellcheck\
       python3-pip python3-setuptools \
    && apt-get autoremove -y  \
    && apt-get clean -y \
    && rm -r /var/cache/* /var/lib/apt/lists/*

# go get to install packages
FROM os-builder as go-builder
RUN set -x \
    && go get -v golang.org/x/tools/gopls 2>&1 \
    && go get -v github.com/uudashr/gopkgs/v2/cmd/gopkgs 2>&1 \
    && go get -v github.com/cweill/gotests/gotests 2>&1 \
    && go get -v github.com/josharian/impl 2>&1 \
    && go get -v honnef.co/go/tools/cmd/staticcheck 2>&1 \
    && go get -v gotest.tools/gotestsum 2>&1 \
    && rm -rf $GOPATH/pkg/* $GOPATH/src/* $HOME/.cache/go-build

# Install global node packages.
FROM go-builder as toolz-builder
RUN set -x \
    && su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g husky commitizen commitlint" 2>&1
# taskctl
ENV TASKCTL_URL=https://raw.githubusercontent.com/taskctl/taskctl/master/install.sh
RUN set -x \
    && mkdir /tools \
    && curl -fsSL $TASKCTL_URL  | sh -s -- -b /tools
#github cli
ENV GH_VER=1.14.0

ENV GH_CLI=https://github.com/cli/cli/releases/download/v${GH_VER}/gh_${GH_VER}_linux_amd64.tar.gz
RUN set -ex \
    && wget -q $GH_CLI -O - | tar xz -C ./ --strip-components 1

#Python Packages
RUN set -ex \
    && pip3 --disable-pip-version-check --no-cache-dir install -r /tmp/packages-tmp/requirements.txt

# shellspec, kcov and gh
FROM shellspec/shellspec-debian:kcov as shellspec-kov

# dev-container
FROM toolz-builder
COPY --from=shellspec-kov   /opt/shellspec              /opt/shellspec
COPY --from=shellspec-kov   /usr/local/bin/kcov*        /usr/local/bin/
COPY --from=shellspec-kov   /usr/local/share/doc/kcov   /usr/local/share/doc/kcov

# goss - Quick and Easy server validation
RUN set -ex \
    && curl -fsSL https://goss.rocks/install | sh

# Sentry CLI
RUN set -ex \
    && curl -fsSL https://sentry.io/get-cli/ | bash

ENV PATH="/opt/shellspec:/tools:$PATH"

ARG USERNAME=vscode
COPY dotfiles/.zshrc /home/$USERNAME/

# Remove library scripts for final image
RUN set -ex \
    && apt-get autoremove --assume-yes \
    && apt-get clean --assume-yes \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /tmp/library-scripts \
    && rm -rf /tmp/packages-tmp
