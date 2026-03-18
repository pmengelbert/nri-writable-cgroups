FROM ubuntu
COPY ./10-writable-cgroups /10-writable-cgroups
COPY ./setup-setup.sh /setup-setup.sh
COPY ./setup-cgroup.sh /setup-cgroup.sh
ENTRYPOINT ["/setup-setup.sh"]
