# Simple usage with a mounted data directory:
# > docker build -t onex .
# > docker run -it -v ~/.onex:/onex/.onex onomy/onex-dev init onex --home /onex/.onex
# Copy genesis.json from dev/config to ~/.onex/config and Dealer and Validator keys are in dev/config
# > docker run -it -v ~/.onex:/onex/.onex onomy/onex-dev keys add dealer --recover --home /onex/.onex
# > docker run -it -v ~/.onex:/onex/.onex onomy/onex-dev keys add validator --recover --home /onex/.onex
# > docker run -it -v ~/.onex:/onex/.onex onomy/onex-dev gentx validator 10000000000000000000stake --chain-id onex --home /onex/.onex
# > docker run -it -v ~/.onex:/onex/.onex onomy/onex-dev collect-gentxs --home /onex/.onex
# > docker run -it -p 26656:26656 -p 26657:26657 -p 1317:1317 -p 9090:9090 -p 9091:9091 -d -v ~/.onex:/onex/.onex onomy/onex-dev start --home /onex/.onex
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /go/src/github.com/onomyprotocol/onex

# Add source files
COPY . .
RUN pwd
RUN ls

RUN go version

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES
RUN make install

# Final image
FROM alpine:edge

ENV onex /onex

# Install ca-certificates
RUN apk add --update ca-certificates

WORKDIR $onex

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/onexd /usr/bin/onexd

EXPOSE 26656
EXPOSE 26657
EXPOSE 1317
EXPOSE 9090
EXPOSE 9091

# Run onexd by default, omit entrypoint to ease using container with onexcli
ENTRYPOINT ["onexd"]