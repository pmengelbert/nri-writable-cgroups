#!/usr/bin/env bash
set -exu

mkdir -p /rootpath/opt/nri/plugins
if pgrep writable; then
    pkill writable
fi
cp -f /10-writable-cgroups /rootpath/opt/nri/plugins
cp -f /setup-cgroup.sh /rootpath/setup-cgroup.sh
chroot /rootpath /setup-cgroup.sh
