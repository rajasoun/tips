#!/usr/bin/env sh

set -e

function prompt_error_msg() {
    MSG=$1
    echo "$MSG" ; exit 1;
}

function install_tips() {
    LATEST_URL="https://github.com/rajasoun/tips/releases/latest"
    LATEST_EFFECTIVE=$(curl -s -L -o /dev/null ${LATEST_URL} -w '%{url_effective}')
    LATEST=${LATEST_EFFECTIVE##*/}

    if [ -z "$TIPS_VER" ]; then
        TIPS_VER=${TIPS_VER:-$LATEST}
    fi
    if [ -z "$TIPS_VER" ]; then
        message="ERROR: Could not automatically detect latest version, set TIPS_VER env var and re-run"
        prompt_error_msg "$message"
    fi

    TIPS_DST=${TIPS_DST:-/usr/local/bin}
    INSTALL_LOC="${TIPS_DST%/}/tips"
    message="ERROR: Cannot write to $TIPS_DST set TIPS_DST elsewhere or use sudo"
    touch "$INSTALL_LOC" || promp_error $message

    arch=""
    if [ "$(uname -m)" = "x86_64" ]; then
        arch="amd64"
    elif [ "$(uname -m)" = "aarch64" ]; then
        arch="arm"
    else
        arch="386"
    fi

    url="https://github.com/rajasoun/tips/releases/download/$TIPS_VER/tips_$TIPS_VER_linux_$arch.tar.gz"

    echo "Downloading $url"
    curl -L "$url" -o "/tmp/tips_$TIPS_VER"
    tar -xvzf /tmp/tips_$TIPS_VER/tips_$TIPS_VER_linux_$arch.tar.gz
    chmod +rx /tmp/tips_$TIPS_VER/tips
    mv /tmp/tips_$TIPS_VER/tips "$INSTALL_LOC"
    rm -fr tmp/tips_$TIPS_VER

    echo "TIPS $TIPS_VER has been installed to $INSTALL_LOC"
    echo "tips --version"
}

install_tips
