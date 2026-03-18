#!/usr/bin/env bash
set -exu

curl -fsSL 'https://github.com/containerd/containerd/releases/download/v2.2.2/containerd-static-2.2.2-linux-amd64.tar.gz' | tar -zxv -C /usr/

cat > /etc/containerd/config.toml <<'EOF'
[plugins.'io.containerd.nri.v1.nri']
# enable nri support in containerd.
    disable = false
# allow connections from externally launched nri plugins.
    disable_connections = false
# plugin_config_path is the directory to search for plugin-specific configuration.
    plugin_config_path = "/etc/nri/conf.d"
# plugin_path is the directory to search for plugins to launch on startup.
    plugin_path = "/opt/nri/plugins"
# plugin_registration_timeout is the timeout for a plugin to register after connection.
    plugin_registration_timeout = "5s"
# plugin_request_timeout is the timeout for a plugin to handle an event/request.
    plugin_request_timeout = "2s"
# socket_path is the path of the nri socket to create for plugins to connect to.
    socket_path = "/var/run/nri/nri.sock"
EOF

mount -t cgroup2 -o remount,rw,nosuid,nodev,noexec,relatime,nsdelegate cgroup2 /sys/fs/cgroup
systemctl restart containerd
