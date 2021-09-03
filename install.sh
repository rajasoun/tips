#!/usr/bin/env bash

set -e

function prompt_error_msg() {
    MSG=$1
    echo "$MSG"
    exit 1
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

    TIPS_DST=${TIPS_DST:-/tools}
    INSTALL_LOC="${TIPS_DST%/}"
    message="ERROR: Cannot write to $TIPS_DST set TIPS_DST elsewhere or use sudo"
    touch "$INSTALL_LOC" || promp_error "$message"

    arch=""
    if [ "$(uname -m)" = "x86_64" ]; then
        arch="amd64"
    elif [ "$(uname -m)" = "aarch64" ]; then
        arch="arm"
    else
        arch="386"
    fi
    LINUX="_linux_$arch.tar.gz"
    TIPS_ARCHIVE="tips_$TIPS_VER$LINUX"
    url="https://github.com/rajasoun/tips/releases/download/$TIPS_VER/$TIPS_ARCHIVE"

    mkdir -p "/tmp/tips"
    TEMP_ARCHIVE="/tmp/tips/$TIPS_ARCHIVE"
    echo "Downloading $url to $TEMP_ARCHIVE"

    curl -L "$url" -o "$TEMP_ARCHIVE"
    tar -xvzf "$TEMP_ARCHIVE" -C "$INSTALL_LOC" tips
    chmod +x "$INSTALL_LOC/tips"
    rm -fr "/tmp/tips"
    echo "TIPS $TIPS_VER has been installed to $INSTALL_LOC"
}

mkdir -p /tools
install_tips
