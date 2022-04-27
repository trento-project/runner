FROM registry.suse.com/bci/golang:1.16 as go-build
WORKDIR /build
# we add what's needed to download go modules so that dependencies can be cached in a dedicate layer
ADD go.mod go.sum /build/
RUN go mod download
ADD . /build
RUN zypper -n in git-core && make build

FROM registry.suse.com/bci/python:3.9 AS trento-runner
RUN /usr/local/bin/python3 -m venv /venv \
    && /venv/bin/pip install 'ansible~=4.6.0' 'requests~=2.26.0' 'rpm==0.0.2' 'pyparsing~=2.0' \
    && zypper -n ref && zypper -n in --no-recommends openssh \
    && zypper -n clean

ENV PATH="/venv/bin:$PATH"
ENV PYTHONPATH=/venv/lib/python3.9/site-packages

# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

COPY --from=go-build /build/trento-runner /usr/bin/trento-runner
LABEL org.opencontainers.image.source="https://github.com/trento-project/runner"
ENTRYPOINT ["/tini", "--", "/usr/bin/trento-runner"]
