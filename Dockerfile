# We expect this to be passed in as a build parameter, but fall back to the algol60 registry
# Usage: --build-arg REGISTRY_HOST=artifactory.algol60.net/docker.io/library/
ARG REGISTRY_HOST=bikeshack.azurecr.io/docker.io/library/

## Build stage should use a standard golang base container that uses alpine
## as the base image.  This is because we want to use the same base image
## for all of our containers to reduce the attack surface.
## Building in a container also enforces a repeatable build process.

# Build base just has the packages installed we need.
FROM cgr.dev/chainguard/go:1.20 as builder



# We need to know the github location of the code for the build
# Usage: --build-arg GITHUB_REPO=github.com/bikeshack/csm-inventory
ARG GITHUB_REPO=github.com/bikeshack/dcim

# Version information
# Usage: --build-arg VERSION=$(git describe --tags --always --dirty)
ARG VERSION=0.0.0

# Git commit information
# Usage: --build-arg GIT_COMMIT=$(git rev-parse HEAD)
ARG GIT_COMMIT=unknown

# Build Platform Information for container labels
# Usage: --build-arg buildDate=$(date +'%Y-%m-%d')
ARG buildDate

# Establish a go environment
RUN go env -w GO111MODULE=on

## Copy your code into the container
COPY internal $GOPATH/src/${GITHUB_REPO}/internal
COPY *.go $GOPATH/src/${GITHUB_REPO}/
COPY pkg $GOPATH/src/${GITHUB_REPO}/pkg
COPY go.mod $GOPATH/src/${GITHUB_REPO}/
COPY go.sum $GOPATH/src/${GITHUB_REPO}/

# Build the actual application(s)
ENV CGO_ENABLED=0 GOOS=linux ARCH=amd64
RUN set -ex \
    && cd $GOPATH/src/${GITHUB_REPO} \
    && go build -ldflags="-s -w" -tags=containers -o /dcim . 

# Final stage is the actual container we will run
FROM cgr.dev/chainguard/wolfi-base:latest

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

VOLUME /migrations
COPY --from=builder /dcim /bin/dcim
RUN apk update && apk add --no-cache --update-cache tini
# Tini is now available at /sbin/tini
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/bin/dcim", "serve"]


# binaries in our containers should run as nobody 65534:65534
USER 65534:65534
