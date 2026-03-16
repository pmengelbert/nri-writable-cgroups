#!/usr/bin/env bash

pod_json="${1:-./pod.json}"
container_json="${2:-./container.json}"

pod_id="$(sudo crictl runp "$pod_json")"
cid="$(sudo crictl create "$pod_id" "$container_json" "$pod_json")"
sudo crictl start "$cid"
sudo crictl exec -it "$cid" /bin/sh
