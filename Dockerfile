#----------------------#
# Deps, Build and Test #
#----------------------#
FROM debian:buster-slim AS build-and-test

RUN apt-get update        \
    && apt-get install -y \
    golang                \
    build-essential       \
    git                   \
    ca-certificates       \
    curl                  \
    unzip

# Install terraform
# TODO: (rem) use `terraform-bundle`
ENV TERRAFORM_VERSION 0.12.3
RUN curl -sLO https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
RUN unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip
RUN mv terraform /usr/local/bin/

# Install JQ
ENV JQ_VERSION 1.6
RUN curl https://github.com/stedolan/jq/releases/download/jq-${JQ_VERSION}/jq-linux64 -Lo /usr/local/bin/jq
RUN chmod +x /usr/local/bin/jq

## Install YQ
ENV YQ_VERSION 2.7.2
RUN curl https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 -Lo /usr/local/bin/yq
RUN chmod +x /usr/local/bin/yq

## Install Goss
ENV GOSS_VERSION v0.3.7
RUN curl -L https://github.com/aelsabbahy/goss/releases/download/${GOSS_VERSION}/goss-linux-amd64 -o /usr/local/bin/goss
RUN chmod +rx /usr/local/bin/goss

# Setup non-root build user
RUN addgroup --quiet build && adduser --quiet --disabled-password --gecos "" --ingroup build build

# Create golang src directory
RUN mkdir -p /go/src/github.com/controlplaneio/simulator-standalone

# Create an empty public ssh key file for the tests
RUN mkdir -p /home/build/.ssh && echo  "ssh-rsa FOR TESTING" > /home/build/.ssh/id_rsa.pub
# Create module cache and copy manifest files
RUN mkdir -p /home/build/go/pkg/mod
COPY ./go.mod /go/src/github.com/controlplaneio/simulator-standalone
COPY ./go.sum /go/src/github.com/controlplaneio/simulator-standalone

# Give ownership of module cache and src tree to build user
RUN chown -R build:build /go/src/github.com/controlplaneio/simulator-standalone
RUN chown -R build:build /home/build/go

# Run all build and test steps as build user
USER build

# Install golang module dependencies before copying source to cache them in their own layer
WORKDIR /go/src/github.com/controlplaneio/simulator-standalone
RUN go mod download

# Add the full source tree
ADD .  /go/src/github.com/controlplaneio/simulator-standalone/
WORKDIR /go/src/github.com/controlplaneio/simulator-standalone/

# TODO: (rem) why is this owned by root after the earlier chmod?
USER root
RUN chown -R build:build /go/src/github.com/controlplaneio/simulator-standalone/

USER build

# Golang build and test
WORKDIR /go/src/github.com/controlplaneio/simulator-standalone
ENV GO111MODULE=on
RUN make test

#------------------#
# Launch Container #
#------------------#
FROM debian:buster-slim

ENV TERRAFORM_VERSION 0.12.3
ENV JQ_VERSION 1.6
ENV YQ_VERSION 2.7.2
ENV GOSS_VERSION v0.3.7

RUN \
  DEBIAN_FRONTEND=noninteractive \
    apt update && apt install --assume-yes --no-install-recommends \
      bash \
      awscli \
      bzip2 \
      file \
      ca-certificates \
      curl \
      gettext-base \
      golang \
      lsb-release \
      make \
      openssh-client \
      gnupg \
 && rm -rf /var/lib/apt/lists/*

# Use 3rd party dependencies from build
COPY --from=build-and-test /usr/local/bin/jq /usr/local/bin/jq
COPY --from=build-and-test /usr/local/bin/yq /usr/local/bin/yq
COPY --from=build-and-test /usr/local/bin/goss /usr/local/bin/goss
COPY --from=build-and-test /usr/local/bin/terraform /usr/local/bin/terraform

# Copy statically linked simulator binary
COPY --from=build-and-test /go/src/github.com/controlplaneio/simulator-standalone/dist/simulator /usr/local/bin/simulator

# Setup non-root launch user
RUN useradd -ms /bin/bash launch
RUN mkdir /app
RUN chown -R launch:launch /app

# Add terraform and perturb/scenario scripts to the image
COPY ./Makefile /app/Makefile
ADD ./terraform /app/terraform
ADD ./simulation-scripts /app/simulation-scripts

# Add goss.yaml to verify the container
ADD ./goss.yaml /app

# Add simulator.yaml config file
# TODO: (rem) make this build-time configurable with an ARG
ADD ./simulator.yaml /app

ENV SIMULATOR_MANIFEST_PATH /app/simulation-scripts/
ENV SIMULATOR_TF_DIR /app/terraform
ENV TF_VAR_shared_credentials_file /app/credentials

USER launch
RUN ssh-keygen -f /home/launch/.ssh/id_rsa -t rsa -N ''
WORKDIR /app

CMD [ "/bin/bash" ]
