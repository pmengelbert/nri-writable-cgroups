#!/usr/bin/env bash
set -exu

mkdir -p /rootpath/opt/nri/plugins
cp /10-writable-cgroups /rootpath/opt/nri/plugins
cp /setup-cgroup.sh /rootpath/setup-cgroup.sh
chroot /rootpath /setup-cgroup.sh
