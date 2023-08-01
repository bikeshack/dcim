# This dockerfile is just for gorleaser https://goreleaser.com/errors/docker-build/
# We expect this to be passed in as a build parameter, but fall back to the algol60 registry
# Usage: --build-arg REGISTRY_HOST=artifactory.algol60.net/docker.io/library/
ARG REGISTRY_HOST=artifactory.algol60.net/docker.io/library/

# Final stage is the actual container we will run
FROM ${REGISTRY_HOST}alpine:3.15

# Labels make it easier to troubleshoot the container http://label-schema.org/rc1/
LABEL schema-version org.label-schema.schema-version="1.0" 
LABEL maintainer="Bikeshack Development Team <info@bikeshack.dev>"
LABEL build-date org.label-schema.build-date=$buildDate
LABEL vcs-ref org.label-schema.vcs-ref=$GIT_COMMIT
LABEL vendor org.label-schema.vendor="Bikeshack Industries"
LABEL version org.label-schema.version=$VERSION
# User Servicable Parts Start Here
LABEL name org.label-schema.name="bikeshack/dcim"
LABEL description org.label-schema.description="CSM Data Center Infrastructure Manager"
LABEL url org.label-schema.url="https://bikeshack.dev/dcim"
LABEL vcs-url org.label-schema.vcs-url="https://github.com/bikeshack/dcim" 

COPY dcim /bin/dcim
VOLUME /migrations
RUN apk add --no-cache tini
# Tini is now available at /sbin/tini
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/bin/dcim", "serve"]


# binaries in our containers should run as nobody 65534:65534
USER 65534:65534