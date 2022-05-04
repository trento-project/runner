#!BuildTag: trento/trento-runner:latest
#!BuildTag: trento/trento-runner:%%VERSION%%
#!BuildTag: trento/trento-runner:%%VERSION%%-build%RELEASE%
FROM bci/bci-base:15.3

# Define labels according to https://en.opensuse.org/Building_derived_containers
# labelprefix=com.suse.trento
PREFIXEDLABEL org.opencontainers.image.title="Trento runner"
PREFIXEDLABEL org.opencontainers.image.description="This contains trento runner %%VERSION%%"
PREFIXEDLABEL org.opensuse.reference="registry.suse.com/trento/trento-runner:%%VERSION%%"
PREFIXEDLABEL org.openbuildservice.disturl="%DISTURL%"
PREFIXEDLABEL org.opencontainers.image.created="%BUILDTIME%"

# Workaround for https://github.com/openSUSE/obs-build/issues/487
RUN zypper --non-interactive in sles-release
RUN zypper --non-interactive in tar gzip make go

COPY runner-%%VERSION%%.tar.gz premium-checks.tar.gz* /build/
COPY vendor.tar.gz /build/

WORKDIR /build/
RUN tar -xvzf runner-%%VERSION%%.tar.gz
RUN tar -xvzf vendor.tar.gz -C runner-%%VERSION%%/

# The next RUN only works if the premium-checks.tar.gz file is present
RUN tar --transform 'flags=r;s/^premium-checks//' -xvzf premium-checks.tar.gz -C runner-%%VERSION%%/ || true

WORKDIR /build/runner-%%VERSION%%
RUN VERSION=%%VERSION%% FLAVOR=Premium make build
RUN cp /build/runner-%%VERSION%%/trento-runner /usr/bin/trento-runner

# Install Ansible dependecies
RUN zypper --non-interactive in openssh python3 python3-rpm python3-requests python3-PyYAML python3-Jinja2 ansible-core tini

# Temporary solution until we have ansible-4.6.0 packaged in IBS
COPY ansible-4.6.0.tar.gz /tmp/ansible-4.6.0.tar.gz
RUN zypper --non-interactive in python3-pip
RUN pip install --no-index --find-links /tmp /tmp/ansible-4.6.0.tar.gz

LABEL org.opencontainers.image.source="https://github.com/trento-project/runner"
ENTRYPOINT ["/tini", "--", "/usr/bin/trento-runner"]
